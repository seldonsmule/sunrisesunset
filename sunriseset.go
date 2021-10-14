package main

/* simple get sunrise/set info */

import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "os"
        "flag"
        "io/ioutil"
        "encoding/json"
        "encoding/base64"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/restapi"
        "github.com/seldonsmule/securityspy"

)

type RiseSet struct {

  timeRise time.Time
  timeSet time.Time

}

func testLockfile() bool {

  lockfile := fmt.Sprintf("%s/tmp/sun.lck", os.Getenv("HOME"))

  _, statErr := os.Stat(lockfile)

  if(os.IsNotExist(statErr)){
    return false
  }

  return true;

}

func deleteLockfile(){

  lockfile := fmt.Sprintf("%s/tmp/sun.lck", os.Getenv("HOME"))

  _, statErr := os.Stat(lockfile)

  // if lock file already exist - just log it and exit
  if(statErr == nil){

    //fmt.Println("Lockfile ", lockfile, " created: ",info.ModTime())
    os.Remove(lockfile);

    return;

  }else{
    fmt.Println("Lockfile already deleted");
  }

}

func createLockfile(){

  lockfile := fmt.Sprintf("%s/tmp/sun.lck", os.Getenv("HOME"))

  info, statErr := os.Stat(lockfile)

  // if lock file already exist - just log it and exit
  if(statErr == nil){

    fmt.Println("Lockfile ", lockfile, " created: ",info.ModTime())

    return;

  }

  // otherwise create it

  lockWriteFile, openErr := os.Create(lockfile)

  if(openErr != nil){

    fmt.Println("Error creating lockfile: ", lockfile );

    return;

  }

  fmt.Println("Created lockfile: ", lockfile );

  lockWriteFile.Close()


}

func moveCamera(ss *securityspy.SecuritySpy, nCameraNum int, nPresetNum int){

    ss.PresetPTZ(nCameraNum, nPresetNum)
    fmt.Printf("moving camera[%d] to preset[%d]\n", nCameraNum, nPresetNum)

/*
  if(bSunRise){
    //ss.PresetPTZ(3, 1)
    ss.PresetPTZ(nCameraNum, nPresetNum)
    fmt.Println("Sun rise - moving camera")

  }else{
    //ss.PresetPTZ(3, 6)
    ss.PresetPTZ(nCameraNum, nPresetNum)
    fmt.Println("Sun set - moving camera")
  }
*/

}

func getSunTimes() (time.Time, time.Time) {

  var astroMap map[string]interface{}

  jsonfile := fmt.Sprintf("%s/tmp/sun.json", os.Getenv("HOME"))

  info, statErr := os.Stat(jsonfile)

  if(statErr == nil){

    logmsg.Print(logmsg.Info, "Last ", jsonfile, " write - ",info.ModTime())

    today := time.Now()

    if(today.Day() != info.ModTime().Day()){
      logmsg.Print(logmsg.Info, "not today - delete sun.json")
      os.Remove(jsonfile)
    }

  }

  jsonReadFile, openErr := os.Open(jsonfile)

  if(openErr == nil){

    logmsg.Print(logmsg.Debug01, "found the file");

    byteValue, _ := ioutil.ReadAll(jsonReadFile)

    json.Unmarshal([]byte(byteValue), &astroMap)

    jsonReadFile.Close()

    logmsg.Print(logmsg.Debug02, astroMap["sunset"])

  }else{

   logmsg.Print(logmsg.Info,"Need new cache file: ", jsonfile)

   r := restapi.NewGet("sunriseset", "https://weather.cit.api.here.com/weather/1.0/report.json?product=forecast_astronomy&name=DC&app_id=DemoAppId01082013GAL&app_code=AJKnXv84fjrb0KIHawS0Tg")


//  r.DebugOn()

  r.HasInnerMap("astronomy")

  if(r.Send()){

 //   r.Dump()

  }


  logmsg.Print(logmsg.Debug02, "astronomy: ", r.GetValue("astronomy"))


// here is the deal - this stupid thing is an array of maps, not
// something considered in my original Telsa use case.  here is how to
// get the data...

// 1. Get the info as an array

  astroArray := restapi.CastArray(r.GetValue("astronomy"))


  logmsg.Print(logmsg.Debug02, "ArrayLength: ", len(astroArray))


// -- used this for look to find the array we are looking for. uncommit to see
/*
  for k, v := range astroArray {
    fmt.Println(k, "=", v)

  } // end for loop
*/

// 2. Get the map desired


//  fmt.Printf("--------------------------\n")
  astroMap = restapi.CastMap(astroArray[0])

  jsonData, _ := json.Marshal(astroMap)

  //fmt.Println(string(jsonData))

  jsonWriteFile, err := os.Create(jsonfile)

  if err != nil {
     panic(err)
  }
  defer jsonWriteFile.Close()

  jsonWriteFile.Write(jsonData)
  jsonWriteFile.Close()

  } // end else json file read


// used this to figure out the names to extract the sunset/rise info - unommit to see

/*
  for k, v := range astroMap {
    fmt.Println(k, "=", v)
  } // end for loop
*/



// 3. get teh value in the map
/*
  fmt.Printf("--------------------------\n")
  fmt.Printf("sunset[%s]\n", astroMap["sunset"])
  fmt.Printf("sunrise[%s]\n", astroMap["sunrise"])
*/

/*
  sunrise :=  restapi.CastString(astroMap["sunrise"])
  sunset  :=  restapi.CastString(astroMap["sunset"])
*/
  sunrise :=  astroMap["sunrise"].(string)
  sunset  :=  astroMap["sunset"].(string)

  now := time.Now()


  sunrise = strings.TrimSuffix(sunrise,"AM")
  sunset = strings.TrimSuffix(sunset,"PM")
  //fmt.Println("Rise:", sunrise)
  riseArray := strings.Split(sunrise, ":")
  riseHour, _ := strconv.Atoi(riseArray[0])
  riseMin, _ := strconv.Atoi(riseArray[1])
  //fmt.Println(riseHour)
  //fmt.Println(riseMin)
  //fmt.Println("Set:", sunset)
  setArray := strings.Split(sunset, ":")
  setHour, _ := strconv.Atoi(setArray[0])
  setHour += 12
  setMin, _ := strconv.Atoi(setArray[1])
  //fmt.Println(setHour)
  //fmt.Println(setMin)

//     func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time


  timeRise := time.Date(now.Year(),
                    now.Month(),
                    now.Day(),
                    riseHour,
                    riseMin,
                    now.Second(),
                    now.Nanosecond(),
                    now.Location())

  timeSet := time.Date(now.Year(),
                    now.Month(),
                    now.Day(),
                    setHour,
                    setMin,
                    now.Second(),
                    now.Nanosecond(),
                    now.Location())

/*
  fmt.Println("Now:", now)
  fmt.Println("Now:", now.Unix())
  fmt.Println("Sunrise:", timeRise)
  fmt.Println("Sunrise:", timeRise.Unix())
  fmt.Println("Sunset:", timeSet)
  fmt.Println("Sunset:", timeSet.Unix())
*/

  return timeRise, timeSet
}

func help(){

  fmt.Println("sunrisesunset allows you to control a camera that is under")
  fmt.Println("SecuritySpy control via its web interface.  Before using you will need to")
  fmt.Println("do the following:")
  fmt.Printf("\n")
  fmt.Printf("\t1. Use the buildconfig cmd to set up a config file. Default location is in $HOME/tmp\n")
  fmt.Printf("\n")
  fmt.Printf("\t2. Go into SecuritySpy and get the camera number and the preset PTZ numbers\n")
  fmt.Printf("\n")
  fmt.Printf("\t3. Use the daynight cmd (in cron is best) to move the cameras based on time of day\n")
  fmt.Printf("\n")

  fmt.Println("Commands and options:\n")
  fmt.Println("-cmd buildconfig -url urlname -idandpass userid:password [-conffile path/name]")
  fmt.Println("\tBuilds the conf file")
  fmt.Println("-cmd show");
  fmt.Println("\tShow contents of config file")
  fmt.Println("-cmd movecamera -camera num -preset num ")
  fmt.Println("\tMoves a camera to a PTZ preset")
  fmt.Println("-cmd daynight -camera num -presetday num -presetnight num")
  fmt.Println("\tDepending on time of day moves camera between PTZ 2 presets")
  fmt.Println("-cmd lock");
  fmt.Println("\tCreates lockfile to disable time logic");
  fmt.Println("-cmd unlock");
  fmt.Println("\tDeletes lockfile to renable time logic");
  fmt.Println("-cmd help");
  fmt.Println("\tDisplay this help message");
  

}

func buildAuthToken(authString string){

  encodeFileName := fmt.Sprintf("%s/tmp/.sun.token", os.Getenv("HOME"))

  fmt.Println(authString)
  sEnc := base64.StdEncoding.EncodeToString([]byte(authString))
  fmt.Println(sEnc)

  encodeWriteFile, err := os.Create(encodeFileName)

  if err != nil {
     panic(err)
  }
  defer encodeWriteFile.Close()

  encodeWriteFile.Write([]byte(sEnc))
  encodeWriteFile.Close()


}

func getToken(decode bool) string {

  encodeFileName := fmt.Sprintf("%s/tmp/.sun.token", os.Getenv("HOME"))

  encodeReadFile, openErr := os.Open(encodeFileName)

  if(openErr == nil){

    byteValue, _ := ioutil.ReadAll(encodeReadFile)

    stringValue := string(byteValue)

    encodeReadFile.Close()

    if(decode){
      fmt.Println("decode it")
      sDec, _ := base64.StdEncoding.DecodeString(stringValue)
      stringValue = string(sDec)
    }

    return stringValue

  }

  return "notset"

}

func moveDayNight(conf string, encryptkey string, camera int, daypreset int, nightpreset int) {

  if(testLockfile()){

    fmt.Println("Stopping logic - lockfile exist - use unlock to remove");
    logmsg.Print(logmsg.Info, "Stopping logic - lockfile exist - use unlock to remove");

    return

  }

  ss := securityspy.NewEncrypt(conf, encryptkey)

  if(ss == nil){
    fmt.Println("Config file not created - use buildconfig")
    return
  }

  now := time.Now()

  timeRise, timeSet := getSunTimes()

  logmsg.Print(logmsg.Info, "Sunrise: ", timeRise)
  logmsg.Print(logmsg.Info, "Sunset: ", timeSet)

  // test for after sunrise
  if(now.After(timeRise) && now.Before(timeSet)){

    logmsg.Print(logmsg.Info, "Sun has rose")

    moveCamera(ss, camera, daypreset)
  }else{

    moveCamera(ss, camera, nightpreset)
    logmsg.Print(logmsg.Info, "Sun has set")

  }

  return

}

func main() {

  logfile := fmt.Sprintf("%s/tmp/sun.log", os.Getenv("HOME"))
  configfile := fmt.Sprintf("%s/tmp/.sun.conf", os.Getenv("HOME"))
  encryptkey := "1234567890AbcDeF"

  cmdPtr := flag.String("cmd", "help", "Command to run")
  idandpassPtr := flag.String("idandpass", "notset", "SecuritySpy web userid:password")
  urlPtr := flag.String("url", "notset", "url of SecuritySpy webserver")
  confPtr := flag.String("conffile", configfile, "path and name of configfile")

  cameraPtr := flag.Int("camera", 0, "SecuritySpy camera number")
  presetPtr := flag.Int("preset", 0, "SecuritySpy preset number")

  daypresetPtr := flag.Int("daypreset", 2, "SecuritySpy preset number for day")
  nightpresetPtr := flag.Int("nightpreset", 1, "SecuritySpy preset number for night")

  flag.Parse()

  if(*cmdPtr == "help"){
    help()
    os.Exit(1)
  }


  logmsg.SetLogLevel(logmsg.Debug03)

  logmsg.SetLogFile(logfile)

  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "idandpassPtr = ", *idandpassPtr)
  logmsg.Print(logmsg.Info, "urlPtr = ", *urlPtr)
  logmsg.Print(logmsg.Info, "confPtr = ", *confPtr)
  logmsg.Print(logmsg.Info, "cameraPtr = ", *cameraPtr)
  logmsg.Print(logmsg.Info, "presetPtr = ", *presetPtr)
  logmsg.Print(logmsg.Info, "daypresetPtr = ", *daypresetPtr)
  logmsg.Print(logmsg.Info, "nightpresetPtr = ", *nightpresetPtr)
  logmsg.Print(logmsg.Info, "tail = ", flag.Args())

  switch *cmdPtr {

    case "show":
      ss := securityspy.NewEncrypt(*confPtr, encryptkey)
      if( ss != nil){
        ss.DumpConfig()
      }else{
        fmt.Println("Missing config file - use buildconfig")
      }

    case "buildconfig":
      fmt.Println("Build config")

      if(*urlPtr == "notset"){
        fmt.Println("Err:  Missing url paramater")
        os.Exit(2)
      }

      if(*idandpassPtr == "notset"){
        fmt.Println("Err:  Missing idandpass paramater")
        os.Exit(2)
      }

      ss := securityspy.NewBuildConfigEncrypt(*urlPtr, *idandpassPtr, *confPtr,
                                              encryptkey)
      ss.DumpConfig()

    case "movecamera":
      fmt.Println("Move Camera")

      if(*cameraPtr == 0){
        fmt.Println("Err:  Missing camera paramater")
        os.Exit(2)
      }

      if(*presetPtr == 0){
        fmt.Println("Err:  Missing preset paramater")
        os.Exit(2)
      }

      ss := securityspy.NewEncrypt(*confPtr, encryptkey)
      if( ss != nil){
        moveCamera(ss, *cameraPtr, *presetPtr)
      }else{
        fmt.Println("Missing config file - use buildconfig")
      }

    case "daynight":
      fmt.Println("Day/night test move logic")

      if(*cameraPtr == 0){
        fmt.Println("Err:  Missing camera paramater")
        os.Exit(2)
      }

      moveDayNight(*confPtr, encryptkey, *cameraPtr, *daypresetPtr, 
                   *nightpresetPtr)

    case "lock":
      fmt.Println("Locking time logic - use unlock to restore");
      createLockfile();

    case "unlock":
      fmt.Println("Removing lock time logic - going back to normal operations");
      deleteLockfile();

    case "help":
      help()
      os.Exit(0)

    default:
      help()
      os.Exit(2)

  } // end switch


}
