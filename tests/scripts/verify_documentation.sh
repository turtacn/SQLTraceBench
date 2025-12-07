#!/bin/bash
echo "Verifying documentation..."

# Check files exist
FILES=("docs/type_mapping_guide.md" "docs/precision_handling.md" "configs/type_mapping_rules.yaml" "configs/precision_policy.yaml")

for file in "${FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo "Missing file: $file"
        exit 1
    fi
done

echo "Documentation verified."
