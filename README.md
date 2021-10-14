# sunrisesunset
This tool allows the user to control a camera that is being managed by [Bensoft's SecuritySpy software](https://bensoftware.com/securityspy/) via its web interface.  

It is a simple tool that's main function is to move a PTZ camera between two preset PTZ points for a day look and a night look

##Before using you will need to do the following:

	1. Use the buildconfig cmd to set up a config file. Default location is in $HOME/tmp

	2. Go into SecuritySpy and get the camera number and the preset PTZ numbers

	3. Use the daynight cmd (in cron is best) to move the cameras based on time of day

##Commands and options:

-cmd buildconfig -url urlname -idandpass userid:password [-conffile path/name]
	Builds the conf file
-cmd movecamera -camera num -preset num 
	Moves a camera to a PTZ preset
-cmd daynight -camera num -presetday num -presetnight num
	Depending on time of day moves camera between PTZ 2 presets
-cmd lock
	Creates lockfile to disable time logic
-cmd unlock
	Deletes lockfile to renable time logic
-cmd help
	Display this help message
