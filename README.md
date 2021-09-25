# Serverless monitoring of Steam dedicated server lists

This is the source code for the blog post https://www.maygames.net/2021/09/25/serverless-monitoring-of-steam-server-lists/

Build instructions:

1. Edit "template.yaml" Environment - Variables: Enter your Steam Web API Key (REPLACE_WITH_YOUR_API_KEY) and the Steam APP Id (REPLACE_WITH_YOUR_APPID) of the game you wish to monitor.
2. Install AWS SAM CLI, see https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
3. Build using SAM:

sam build

Deploy on AWS using SAM:

sam deploy --guided

The Lambda function is deployed as "SteamMetricsFunction"