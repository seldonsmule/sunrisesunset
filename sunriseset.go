package main

/* simple get sunrise/set info */

import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "os"
        "bytes"
        "flag"
        "io/ioutil"
        "encoding/json"
        "encoding/base64"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/restapi"
        "github.com/seldonsmule/securityspy"

)

type RiseSet struct {

  LocationName string `json:"location_name"`

  Sunrise int `json:"sunrise"`
  Sunset int `json:"sunset"`

}

type WeatherConf struct { // info for using the developer.here.com API

  ApiKey  string
  ZipCode string
  Url string

  confFile string
  dataFound bool

  WeatherFlowUrl string
  WeatherFlowStationId string
  WeatherFlowToken string

}

func NewWeather() *WeatherConf{

  w := new(WeatherConf)

  //w.confFile = fmt.Sprintf("%s/tmp/.developer.here.com.json", os.Getenv("HOME"))
  w.confFile = fmt.Sprintf("%s/tmp/.sunriseset.json", os.Getenv("HOME"))

  w.dataFound = false

  return w

}

func (pW *WeatherConf) SetWeatherFlow(stationid string, token string) {

  pW.WeatherFlowStationId = stationid
  pW.WeatherFlowToken = token
  pW.WeatherFlowUrl = "https://swd.weatherflow.com/swd/rest/better_forecast?station_id=" + stationid + "&token=" + token

  pW.dataFound = true
}

func (pW *WeatherConf) SetWeather(apikey string, zipcode string) {

  pW.ApiKey = apikey
  pW.ZipCode = zipcode
  pW.Url = "https://weather.cit.api.here.com/weather/1.0/report.json?product=forecast_astronomy"

  pW.dataFound = true
}

func (pW *WeatherConf) CompleteUrl() string{

  return(fmt.Sprintf("%s&zipcode=%s&apikey=%s", pW.Url, pW.ZipCode,
                                                 pW.ApiKey))

}

func (pW *WeatherConf) ReadConf() bool{

  pW.dataFound = false

  file, err := os.Open(pW.confFile) // for read access

  if (err != nil){
    logmsg.Print(logmsg.Error,"Unable to open configfile: ", err," ", pW.confFile)

     //fmt.Println("Unable to read confifile: ", err, " ", pW.confFile)
    return false
  }

  defer file.Close()

  data := make([]byte, 1000)

  count, err := file.Read(data)

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to read config: ", err, count)
     return false
  }

  err = json.NewDecoder(bytes.NewReader(data)).Decode(pW)

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to decode config: ", err)
     return false
  }


  pW.dataFound = true

  return true
}

func (pW *WeatherConf) SaveConf() bool{

  j, err := json.Marshal(pW)

  if(err != nil){
    fmt.Println(err)
    return false
  }

  fmt.Println("Saving config: ", pW.confFile)

  writeFile, err := os.Create(pW.confFile)

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to write config: ", err)
     fmt.Println("Unable to write config: ", err)
     return false
  }

  defer writeFile.Close()

  writeFile.Write(j)
  writeFile.Close()

  return true

}

func (pW *WeatherConf) Dump(){

  fmt.Println("Dumping Weather")

  fmt.Println("ApiKey: ", pW.ApiKey)
  fmt.Println("ZipCode: ", pW.ZipCode)
  fmt.Println("Url: ", pW.Url)
  fmt.Println("CompleteUrl: ", pW.CompleteUrl())
  fmt.Println("WeatherFlowStationId: ", pW.WeatherFlowStationId)
  fmt.Println("WeatherFlowToken: ", pW.WeatherFlowToken)
  fmt.Println("WeatherFlowUrl: ", pW.WeatherFlowUrl)


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

func getWeatherFlowSunTimes(pW *WeatherConf) (time.Time, time.Time) {

  var forecast BetterForcast
  var risesettimes RiseSet

  timeRise := time.Now()
  timeSet := time.Now()

  jsonfile := fmt.Sprintf("%s/tmp/weatherflow.json", os.Getenv("HOME"))

  info, statErr := os.Stat(jsonfile)

  // logic so we don't call the API to much 
  // we call it once a day

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
    fmt.Println("jsonReadFile: ", jsonfile)

    byteValue, _ := ioutil.ReadAll(jsonReadFile)

    json.Unmarshal([]byte(byteValue), &risesettimes)

    jsonReadFile.Close()

    logmsg.Print(logmsg.Debug02, "risesettimes.LocationName: ", risesettimes.LocationName)

  }else{

   logmsg.Print(logmsg.Info,"Need new cache file: ", jsonfile)

   // This API comes from weatherflow
   //

   r := restapi.NewGet("sunriseset", pW.WeatherFlowUrl)

   r.JsonOnly()

//  r.DebugOn()

    if(r.Send()){

   //   r.Dump()
   //fmt.Println("r.Response: ", r.GetResponseBody())
  
    }

    json.Unmarshal(r.BodyBytes, &forecast)

    //fmt.Println("forecast: ", forecast)

    fmt.Println("locationName: ", forecast.LocationName)

    risesettimes.LocationName = forecast.LocationName
    risesettimes.Sunrise = forecast.Forecast.Daily[0].Sunrise
    risesettimes.Sunset = forecast.Forecast.Daily[0].Sunset

    jsonData, _ := json.Marshal(risesettimes)

    //fmt.Println(string(jsonData))

    jsonWriteFile, err := os.Create(jsonfile)

    if err != nil {
       panic(err)
    }
    defer jsonWriteFile.Close()

    jsonWriteFile.Write(jsonData)
    jsonWriteFile.Close()

  } //end else

 
  timeRise = time.Unix(int64(risesettimes.Sunrise), 0)
  timeSet = time.Unix(int64(risesettimes.Sunset), 0)


  return timeRise, timeSet
}

func getSunTimes(pW *WeatherConf) (time.Time, time.Time) {

  var astroMap map[string]interface{}

  jsonfile := fmt.Sprintf("%s/tmp/sun.json", os.Getenv("HOME"))

  info, statErr := os.Stat(jsonfile)

  // logic so we don't call the API to much and run up the
  // counter on number of calls beyond the free tier
  // we call it once a day

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

   // This API comes from developer.here.com.  You will need to get your own
   // apikey
   //
   //r := restapi.NewGet("sunriseset", "https://weather.cit.api.here.com/weather/1.0/report.json?product=forecast_astronomy&name=DC&app_id=DemoAppId01082013GAL&app_code=AJKnXv84fjrb0KIHawS0Tg")
//   r := restapi.NewGet("sunriseset", "https://weather.cit.api.here.com/weather/1.0/report.json?product=forecast_astronomy&zipcode=33914&apikey=oXDrTYi4vcE3_yddC5HxAewFzDib4FnBNZRuFbBR2d0")

   r := restapi.NewGet("sunriseset", pW.CompleteUrl())


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
  fmt.Printf("\t2. Use the setweather cmd to set up a weather file. \n")
  fmt.Printf("\n")
  fmt.Printf("\t3. Go into SecuritySpy and get the camera number and the preset PTZ numbers\n")
  fmt.Printf("\n")
  fmt.Printf("\t4. Use the daynight cmd (in cron is best) to move the cameras based on time of day\n")
  fmt.Printf("\n")

  fmt.Println("Commands and options:\n")
  fmt.Println("-cmd buildconfig -url urlname -idandpass userid:password [-conffile path/name]")
  fmt.Println("\tBuilds the camera conf file")
  fmt.Println("-cmd setweather -hereapikey key -zipcode localzip")
  fmt.Println("\tBuilds the weather conf file.  You need an api key from developer.here.com to setup and use")
  fmt.Println("-cmd setweatherflow -weatherstationid stationid -weatherstationtoken token")
  fmt.Println("\tBuilds the weather conf file.  You need the station id and developer token from your WeatherFlow system.  Got to setup to find it.")
  fmt.Println("-cmd show");
  fmt.Println("\tShow contents of config files")
  fmt.Println("-cmd printtimes");
  fmt.Println("\tDisplays sunrise/sunset times")
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

func printTimes(pW *WeatherConf){

  timeRise, timeSet := getSunTimes(pW)

  fmt.Println("Sunrise: ", timeRise)
  fmt.Println("Sunset: ", timeSet)
}

func printWeatherFlowTimes(pW *WeatherConf){

  timeRise, timeSet := getWeatherFlowSunTimes(pW)

  fmt.Println("Sunrise: ", timeRise)
  fmt.Println("Sunset: ", timeSet)
}

func moveDayNight(conf string, encryptkey string, camera int, daypreset int, nightpreset int, pW *WeatherConf) {

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

//  timeRise, timeSet := getSunTimes(pW)
  timeRise, timeSet := getWeatherFlowSunTimes(pW)

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

  fmt.Println("Sunrise/Sunset control for SecuritySpy")

  logfile := fmt.Sprintf("%s/tmp/sun.log", os.Getenv("HOME"))
  configfile := fmt.Sprintf("%s/tmp/.sun.conf", os.Getenv("HOME"))
  encryptkey := "1234567890AbcDeF"

  cmdPtr := flag.String("cmd", "help", "Command to run")
  idandpassPtr := flag.String("idandpass", "notset", "SecuritySpy web userid:password")
  zipPtr := flag.String("zipcode", "notset", "Zipcode for sunrise/sunset")
  hereApiKeyPtr := flag.String("hereapikey", "notset", "developer.here.com apikey used for the weather check")
  urlPtr := flag.String("url", "notset", "url of SecuritySpy webserver")
  confPtr := flag.String("conffile", configfile, "path and name of configfile")

  weatherStationIdPtr := flag.String("weatherstationid", "notset", "Weather station id to use for weather (see weatherflow app for info")
  weatherStationTokenPtr := flag.String("weatherstationtoken", "notset", "Weather station token to use for weather (see weatherflow app for info")

  cameraPtr := flag.Int("camera", 0, "SecuritySpy camera number")
  presetPtr := flag.Int("preset", 0, "SecuritySpy preset number")

  daypresetPtr := flag.Int("daypreset", 2, "SecuritySpy preset number for day")
  nightpresetPtr := flag.Int("nightpreset", 1, "SecuritySpy preset number for night")

  WConf := NewWeather()
  if(!WConf.ReadConf()){
    fmt.Println("Error reading config file - initalizing instead\n\n")
    WConf.SaveConf()
    if(!WConf.ReadConf()){
      fmt.Println("Error reading config file after initalizing - something is wrong\n\n")
      os.Exit(1)
    }
  }
  /*
  if(!WConf.ReadConf()){
    fmt.Println("Error reading config file\n\n")
    os.Exit(1)
  }
  */

  flag.Parse()

  if(*cmdPtr == "help"){
    help()
    os.Exit(1)
  }


  logmsg.SetLogLevel(logmsg.Debug03)

  logmsg.SetLogFile(logfile)

  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "idandpassPtr = ", *idandpassPtr)
  logmsg.Print(logmsg.Info, "zipPtr = ", *zipPtr)
  logmsg.Print(logmsg.Info, "hereApiKeyPtr = ", *hereApiKeyPtr)
  logmsg.Print(logmsg.Info, "urlPtr = ", *urlPtr)
  logmsg.Print(logmsg.Info, "confPtr = ", *confPtr)
  logmsg.Print(logmsg.Info, "weatherStationIdPtr = ", *weatherStationIdPtr)
  logmsg.Print(logmsg.Info, "weatherStationTokenPtr = ", *weatherStationTokenPtr)
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
	os.Exit(2)
      }

      if(WConf.dataFound){
        WConf.Dump()
      }else{
        fmt.Println("ERROR-> Weather Conf data not found")
      }

    case "printtimes":
      printTimes(WConf)

    case "printweatherflowtimes":
      printWeatherFlowTimes(WConf)

    case "setweatherflow":
      if(*weatherStationIdPtr == "notset"){
        fmt.Println("Err:  Missing weatherstationid paramater")
        os.Exit(2)
      }

      if(*weatherStationTokenPtr == "notset"){
        fmt.Println("Err:  Missing weatherstationtoken paramater")
        os.Exit(2)
      }

      WConf.SetWeatherFlow(*weatherStationIdPtr, *weatherStationTokenPtr)
      WConf.Dump()
      if(!WConf.SaveConf()){
        fmt.Println("Err: Unable to save weather file")
        os.Exit(2)
      }

    case "setweather":
      if(*hereApiKeyPtr == "notset"){
        fmt.Println("Err:  Missing hereapikey paramater")
        os.Exit(2)
      }

      if(*zipPtr == "notset"){
        fmt.Println("Err:  Missing zipcode paramater")
        os.Exit(2)
      }

      WConf.SetWeather(*hereApiKeyPtr, *zipPtr)
      WConf.Dump()
      if(!WConf.SaveConf()){
        fmt.Println("Err: Unable to save weather file")
        os.Exit(2)
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

      if(!WConf.dataFound){
        fmt.Println("Err: Missing Weather Conf data - Use -cmd setweather")
        os.Exit(2)
      }

      fmt.Println("Day/night test move logic")

      if(*cameraPtr == 0){
        fmt.Println("Err:  Missing camera paramater")
        os.Exit(2)
      }

      moveDayNight(*confPtr, encryptkey, *cameraPtr, *daypresetPtr, 
                   *nightpresetPtr, WConf)

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
