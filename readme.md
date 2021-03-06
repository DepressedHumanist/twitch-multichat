## Twitch multichat
### What does it do?
It is a small app that allows you to see chats of multiple streamers in one chat window.

### Why would I need that?
I made it mostly for myself, but it may be useful for small streamers doing collabs, so that 
they can easily read (and respond to) the chats of other streamers in the same collab. 

### How to use it?
First, download it from [here](https://github.com/DepressedHumanist/twitch-multichat/releases).  

Then unzip it to some foldeer and launch the `.exe` file. You will need to authenticate
via twitch, but I won't collect any of your data, it is only required to be able to get 
streamer avatars and emotes. After authentication you will be greeted with a small 
configuration window where you can input nicknames of the streamers that you want to read. 
This window will always appear on startup unless you check the "Do not show this window
on startup" checkbox.  

When you submit the list of streamers, the chat window will appear. If you wish to return 
to the configuration screen, hover your mouse around the top of the window and the settings 
button will appear, click it to go to the configuration page.

### Can I use it with OBS/other streaming software?
Yes. What you see in the app window is basically a webpage, you can access the chat page 
by going into your browser and navigating to `localhost:8081/chat`. You can use this 
URL in browser capture of your favorite streaming software to get transparent background.

###Can I change the style of the chat?
Yes, kinda. You can change the file `assets/chat.css` to get the desired look of your chat,
but I'm not sure how compatible it is with existing stylesheets. Message me or open an 
issue if you really want to do it but have some troubles!

### Features
This app supports the display of not only twitch emotes, but BetterTTV and FrankerFaceZ 
emotes. It will display shared and channel emotes from BTTV/FFZ for every channel in the 
collab, even if one or several of the other streamers don't have this emote enabled.

To quickly see which chat the message came from, streamer's avatar is shown before every 
message to indicate that.

More features to come:
- badges
- filtering bots
- highlighting subscriptions, bits, etc
- easier customization
- and more! You can even request something if you want to see it
