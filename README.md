# Notes

Self-hosted note taking application.

# Development

Run the application stack. Note that development data is note saved to a volume and will not survive host reboots.

```
docker compose up
```

# Deployment

A bespoke script is available at `./deploy.sh` to deploy to a kubernetes cluster. The deployment relies on a k8s secret, `notes-secrets`, existing in the namespace.

# Contemporaries

- https://keep.google.com/
- https://bear.app/
- https://obsidian.md/
- https://app.simplenote.com/
- https://joplinapp.org/
- https://evernote.com/
- https://www.onenote.com
