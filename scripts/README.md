# Scripts directory

This places is becoming a bit of experimentation and test data generate scripts.

Here is one

```
# Install required packages
npm install --save-dev @types/node typescript

# Create minimal tsconfig.json
echo '{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "esModuleInterop": true
  }
}' > tsconfig.json

# Compile
npx tsc

# Run the compiled JavaScript
node dist/generate-test-data.js
```