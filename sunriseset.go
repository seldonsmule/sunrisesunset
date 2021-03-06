package main

/* simple get sunrise/set info */

import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "os"
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

func moveCamera(ss *securityspy.SecuritySpy, bSunRise bool){

  if(bSunRise){
    ss.PresetPTZ(3, 1)
    fmt.Println("Sun rise - moving camera to street")

  }else{
    ss.PresetPTZ(3, 6)
    fmt.Println("Sun set - moving camera to garage")
  }

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

  fmt.Println("sunsetrise [options] [auth userid:password] | [show] | [day] | [\n")
  fmt.Println("With no parms - execute sunrise/set test and camera move\n")
  fmt.Println("auth userid:password - Base64 encodes and saves userid and password\n");
  fmt.Println("buildconfig url userid:password")
  fmt.Println("day - Move camera to day possition\n");
  fmt.Println("night - Move camera to night possition\n");
  fmt.Println("lock - Creates lockfile to disable time logic\n");
  fmt.Println("unlock - Deletes lockfile to renable time logic\n");
  fmt.Println("help - Display this\n");
  

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

func main() {

  logfile := fmt.Sprintf("%s/tmp/sun.log", os.Getenv("HOME"))
  configfile := fmt.Sprintf("%s/tmp/.sun.conf", os.Getenv("HOME"))
  encryptkey := "1234567890AbcDeF"

  logmsg.SetLogLevel(logmsg.Debug03)

  logmsg.SetLogFile(logfile)

  args := os.Args;

  if(len(args) >=2){ // being used for other reasons than moving the camera

    switch args[1]{

      case "lock":
        fmt.Println("Locking time logic - use unlock to restore");
        createLockfile();

      case "unlock":
        fmt.Println("Removing lock time logic - going back to normal operations");
        deleteLockfile();

      case "show":
        ss := securityspy.NewEncrypt(configfile, encryptkey)
        if( ss != nil){
          ss.DumpConfig()
        }else{
          fmt.Println("Missing config file - use buildconfig")
        }

      case "buildconfig":
        if(len(args) == 4){
          fmt.Println("Build config")
          ss := securityspy.NewBuildConfigEncrypt(args[2], args[3], configfile,
                                                encryptkey)
          ss.DumpConfig()
        }else{
          help()
        }

      case "day":
        fmt.Println("day move")
        ss := securityspy.NewEncrypt(configfile, encryptkey)
        if( ss != nil){
          moveCamera(ss, true)
        }else{
          fmt.Println("Missing config file - use buildconfig")
        }

      case "night":
        fmt.Println("night move")
        ss := securityspy.NewEncrypt(configfile, encryptkey)
        if( ss != nil){
          moveCamera(ss, false)
        }else{
          fmt.Println("Missing config file - use buildconfig")
        }

      default:
        help()

    }

    os.Exit(0)

  }

  if(testLockfile()){

    fmt.Println("Stopping logic - lockfile exist - use unlock to remove");
    logmsg.Print(logmsg.Info, "Stopping logic - lockfile exist - use unlock to remove");

    os.Exit(0);

  }

  ss := securityspy.NewEncrypt(configfile, encryptkey)

  if(ss == nil){
    fmt.Println("Config file not created - use buildconfig")
    os.Exit(0)
  }

  now := time.Now()

  timeRise, timeSet := getSunTimes()

  logmsg.Print(logmsg.Info, "Sunrise: ", timeRise)
  logmsg.Print(logmsg.Info, "Sunset: ", timeSet)

  // test for after sunrise
  if(now.After(timeRise) && now.Before(timeSet)){

    logmsg.Print(logmsg.Info, "Sun has rose")

   moveCamera(ss,true)
  }else{

    moveCamera(ss,false)
    logmsg.Print(logmsg.Info, "Sun has set")

  }



}
