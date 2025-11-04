Given all of my lessons learned (my original uploads) as well as your screener files, can you help me architect a new Go / Fyne app for Windows?  I've attached a screenshot of my current (failing) project.  Please create a architectural-overview.md file.  At the top I need there to be a section for rules, such as "do not create new extra features if not explicitly asked for or at least approved by the user architect"



Desired features: 

- sleek styling
- day and night mode, each with green themes.  maybe forest / british racing green as the dark hue for night mode and something lighter for day mode, but the emphasis should be on textual letters, numbers, characters contrasting well against the backdrop of whatever screen, button, or popup they're on
- a workflow that goes from screen to screen (or tab to tab): 
  - starts off with what would you like to trade (as in sector)
  - leads to buttons for our matching screeners that perform well for that sector (pointing to finviz in chart mode where i'm immediately greeted by charts)
    - here's one example that opens the charts automatically, although I don't know - v=211&p=d&s=ta_newlow&f=cap_largeover,sh_avgvol_o1000,sh_price_o20,ta_sma50_pb,ta_sma200_pb&o=-relativevolume
  - goes to a screen where we enter the ticker symbol and select from a dropdown the matching strategy / pine script.  ideally these options would be correlated or dictated by the sector we chose two steps ago so that we can only pick pine scripts that marry up well with the sector (and security) we're trading 
    - start cooldown timer to limit impulsivity
  - flows into a checklist that will help to hinder impulsivity that echoes trend-following principles
    - ie (see attached screenshots)
  - goes into a position sizing tab that helps to calculate position size based on tried and true poker betting principles, but for options trading though
  - next up is a heat check screen that will ensure we're not too heavy in one area.  This won't be pure basket trading though, as our backtesting has demonstrated certain sectors are blackout zones
  - then our trade entry screen, where we select the options strategy from: 
    - bull call spread
    - bear put spread
    - bull put credit spread
    - bear call credit spread
    - long call
    - long put
    - covered call
    - cash-secured put
    - iron butterfly
    - iron condor
    - long put butterfly
    - long call butterfly
    - calendar call spread
    - calendar put spread
    - diagonal call spread
    - diagonal put spread
    - inverse iron butterfly
    - inverse iron condor 
    - short put butterfly
    - short call butterfly
    - straddle
    - strangle
    - call ratio backspread
    - put broken wing
    - put ratio backspread
    - call broken wing
  - lastly and most importantly a calendar screen, where we can see horserace style lines for each trade across sectors.  For example, it would be based on the Windows computer's current date and time.  It would show two weeks back in time, and 12 weeks forward in time.  In that way we will have a grid where the sectors are the y axis (up and down) and the ticker symbols that we enter are the x axis (left to right).  So we might have a butterfly trade in tech that started last week that's expiring in two weeks, a bull call spread trade in healthcare that started this week and which goes for four weeks, and so on.  In this way, we can see at a glance our current position and diversification
  - the real final screen (although we shouldn't need it that often) should be where we can edit or delete past trades.  It is for that reason that the trades should be saved automatically as we progress across the screens, so we don't lose our progress
- I also want a few other features - sample data (sample trades) at the push of a button
  - vimium mode (like the browser extension by phil crosby) and toggle button for on/off
  - a help button with a question mark
  - a welcome screen that is shown on startup but which also has a corresponding button in case you miss it

