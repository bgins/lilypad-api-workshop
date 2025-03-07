# lilypad-api-workshop

This project demonstrates a simple call to the Lilypad Anura API.

### Usage

Visit https://anura.lilypad.tech/ to request an API key.

Set your API key and the model you would like to run as environment variables.


```bash
# Use a specific model
export ANURA_MODEL="gpt-4"
export ANURA_API_KEY="<your-api-key>"
```

We default to `qwen2.5:7b`. See the Anura docs to see what other models are available: https://docs.lilypad.tech/lilypad/developer-resources/inference-api#get-available-models.

Run a job:
```
go run .
```

### Credits

Most of this code was written by Claude 3.7 Sonnet. I haven't looked it over that closely, so take it with a grain of salt.
