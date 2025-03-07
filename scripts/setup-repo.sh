#!/bin/bash

set -e

# Default values
GIT_REPO=""
GIT_BRANCH="main"
CONFIG_DIR="configs"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --repo=*)
      GIT_REPO="${1#*=}"
      shift
      ;;
    --branch=*)
      GIT_BRANCH="${1#*=}"
      shift
      ;;
    --config-dir=*)
      CONFIG_DIR="${1#*=}"
      shift
      ;;
    --help)
      echo "Usage: $0 --repo=<git-repo-url> [--branch=<branch>] [--config-dir=<dir>]"
      echo
      echo "Options:"
      echo "  --repo=<url>        Git repository URL (required)"
      echo "  --branch=<branch>   Git branch to use (default: main)"
      echo "  --config-dir=<dir>  Directory for configurations (default: configs)"
      echo "  --help              Display this help message"
      exit 0
      ;;
    *)
      echo "Error: Unknown option $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

# Check if Git repository URL is provided
if [[ -z "$GIT_REPO" ]]; then
  echo "Error: Git repository URL is required"
  echo "Use --help for usage information"
  exit 1
fi

echo "ConfigSync Repository Setup"
echo "=========================="
echo
echo "Git Repository: $GIT_REPO"
echo "Git Branch: $GIT_BRANCH"
echo "Config Directory: $CONFIG_DIR"
echo

# Check if Git is installed
if ! command -v git &> /dev/null; then
  echo "Error: git is not installed or not in PATH"
  exit 1
fi

# Create a ConfigSync secret with Git credentials (if needed)
echo "Do you need to configure Git credentials? [y/N]"
read -r NEED_CREDENTIALS

if [[ "$NEED_CREDENTIALS" =~ ^[Yy]$ ]]; then
  echo "Enter Git username:"
  read -r GIT_USERNAME
  
  echo "Enter Git password/token:"
  read -rs GIT_PASSWORD
  echo
  
  # Create Kubernetes secret with Git credentials
  echo "Creating Git credentials secret..."
  kubectl create secret generic configsync-git-credentials \
    --namespace=configsync-system \
    --from-literal=username="$GIT_USERNAME" \
    --from-literal=password="$GIT_PASSWORD" \
    --dry-run=client -o yaml | kubectl apply -f -
  
  echo "Git credentials secret created."
fi

# Create example ConfigSync resource
echo "Creating example ConfigSync resource..."
cat <<EOF > "$(dirname "$0")/../examples/configsync-example.yaml"
apiVersion: configsync.io/v1alpha1
kind: ConfigSync
metadata:
  name: example-configsync
  namespace: default
spec:
  gitRepository: "${GIT_REPO}"
  branch: "${GIT_BRANCH}"
  path: "${CONFIG_DIR}"
  environment: "development"
  syncInterval: "5m"
  autoApprove: false
  driftDetection: true
EOF

echo "Example ConfigSync resource created at examples/configsync-example.yaml"
echo

# Check if the user wants to create an initial directory structure
echo "Do you want to create an initial directory structure in your Git repository? [y/N]"
read -r CREATE_STRUCTURE

if [[ "$CREATE_STRUCTURE" =~ ^[Yy]$ ]]; then
  TEMP_DIR=$(mktemp -d)
  
  echo "Cloning repository to create initial structure..."
  if ! git clone "$GIT_REPO" -b "$GIT_BRANCH" "$TEMP_DIR"; then
    echo "Error: Failed to clone repository"
    rm -rf "$TEMP_DIR"
    exit 1
  fi
  
  # Create directory structure
  mkdir -p "$TEMP_DIR/$CONFIG_DIR"/{dev,staging,prod}
  
  # Create README files
  cat <<EOF > "$TEMP_DIR/$CONFIG_DIR/README.md"
# ConfigSync Configurations

This directory contains Kubernetes configurations managed by ConfigSync.

## Directory Structure

- \`dev/\`: Development environment configurations
- \`staging/\`: Staging environment configurations
- \`prod/\`: Production environment configurations

## Usage

Add your Kubernetes manifests to the appropriate environment directory.
ConfigSync will automatically apply these configurations to your cluster.
EOF

  # Create example configuration
  cat <<EOF > "$TEMP_DIR/$CONFIG_DIR/dev/example-configmap.yaml"
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-config
  namespace: default
data:
  example.key: "example value"
  environment: "development"
EOF

  # Commit and push changes
  cd "$TEMP_DIR"
  git add .
  git commit -m "Initialize ConfigSync directory structure"
  
  echo "Pushing changes to repository..."
  if ! git push; then
    echo "Error: Failed to push changes to repository"
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
    exit 1
  fi
  
  cd - > /dev/null
  rm -rf "$TEMP_DIR"
  echo "Initial directory structure created and pushed to repository."
fi

echo
echo "ConfigSync repository setup complete!"
echo
echo "Next steps:"
echo "1. Apply the example ConfigSync resource:"
echo "   kubectl apply -f examples/configsync-example.yaml"
echo
echo "2. Check the ConfigSync status:"
echo "   kubectl get configsyncs -n default"
echo

exit 0 