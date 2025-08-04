# Google Authentication Setup

## Using Application Default Credentials (ADC)

1. Install Google Cloud CLI:
   ```bash
   # macOS
   brew install google-cloud-sdk
   
   # Linux
   curl https://sdk.cloud.google.com | bash
   ```

2. Authenticate with your Google account:
   ```bash
   gcloud auth application-default login
   ```

3. This will open a browser and save credentials to:
   - macOS/Linux: `~/.config/gcloud/application_default_credentials.json`

4. Run the CLI tool:
   ```bash
   ./pkit delete-spam
   ```

## Alternative: Service Account

1. Create a service account in Google Cloud Console
2. Download the JSON key file
3. Set environment variable:
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
   ```