#!/bin/bash

set -e

# Display a welcome message
echo "ConfigSync Installation Script"
echo "============================="
echo

# Check for required tools
echo "Checking for required tools..."
for tool in kubectl helm git; do
  if ! command -v $tool &> /dev/null; then
    echo "Error: $tool is not installed or not in PATH"
    exit 1
  fi
done
echo "All required tools are available."

# Get Kubernetes cluster info
echo "Checking Kubernetes cluster connection..."
if ! kubectl cluster-info &> /dev/null; then
  echo "Error: Could not connect to Kubernetes cluster"
  exit 1
fi
echo "Connected to Kubernetes cluster."

# Create the operator namespace
echo "Creating ConfigSync namespace..."
kubectl create namespace configsync-system --dry-run=client -o yaml | kubectl apply -f -

# Apply CRDs
echo "Installing ConfigSync CRDs..."
kubectl apply -f $(dirname "$0")/../deploy/operator.yaml

# Check if namespace has been created
echo "Verifying installation..."
if ! kubectl get namespace configsync-system &> /dev/null; then
  echo "Error: Failed to create configsync-system namespace"
  exit 1
fi

echo "Creating sample environment directories..."
mkdir -p $(dirname "$0")/../configs/{dev,staging,prod}

# Installation complete
echo
echo "ConfigSync installation complete!"
echo
echo "Next steps:"
echo "1. Configure your Git repository:"
echo "   ./scripts/setup-repo.sh --repo=<your-git-repo-url>"
echo
echo "2. Create a ConfigSync resource:"
echo "   kubectl apply -f examples/configsync-example.yaml"
echo
echo "3. Check the operator status:"
echo "   kubectl -n configsync-system get pods"
echo

exit 0 