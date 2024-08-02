# Generic Twitch Chatbot - GTC

## Available commands

- `!ping` - if the bot is working, just reply with pong!
- `!weather "city"` - displays the weather in the specified city (by default Innopolis)
- `!drawing min:sec` - provide an opportunity to hold a drawing with certain duration
- `!yt` - displays a link to the latest YouTube video from specified channel

## Deployment

During the first launch of the program, you will be greeted by a form in which you will need to enter your Twitch, OAuth token of your bot account, Google API key and MapTiler API key. You can find all necessary links in references.

The most convenient way is by using Docker:

```Shell
git clone git@github.com:setanoier/generic.tv.git
cd generic.tv
docker build -t chatbot .
docker run -td chatbot
```



## References
- Here you can grab an OAuth token for your bot account: https://twitchapps.com/tmi/
- Google API key: https://console.cloud.google.com/apis/credentials
- YouTube Data API v3 (should be enabled): https://console.cloud.google.com/apis/library/youtube.googleapis.com
- MapTiler API key: https://cloud.maptiler.com/account/keys/
