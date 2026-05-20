#!/bin/sh
set -e

CONTACT_POINTS="/etc/grafana/provisioning/alerting/contact-points.yml"

if [ -n "$DISCORD_WEBHOOK_URL" ]; then
    echo "Injecting Discord webhook URL into contact points..."
    sed -i "s|__DISCORD_WEBHOOK_URL__|$DISCORD_WEBHOOK_URL|g" "$CONTACT_POINTS"
    echo "Discord webhook URL injected successfully."
else
    echo "DISCORD_WEBHOOK_URL not set."
    echo "Alerts will fire in Grafana UI but Discord notifications will not be sent."
    echo "Set DISCORD_WEBHOOK_URL environment variable and restart to enable Discord notifications."
fi

exec "/run.sh"
