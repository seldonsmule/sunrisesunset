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

)

type RiseSet struct {

  timeRise time.Time
  timeSet time.Time

}

func moveCamera(bSunRise bool){

   //r := restapi.NewGet("sunriseset", "https://macdaddy.home.c-the-world.org:8001/++ptz/command?cameraNum=3&command=12")
  // r := restapi.NewGet("sunriseset", "https://192.168.2.39:8001/++ptz/command?cameraNum=3&command=17")

  url := "https://macdaddy.home.c-the-world.org:8001/++ptz/command?cameraNum=3&command="


  if(bSunRise){
    fmt.Println("Sun rise - moving camera to street")

    url = fmt.Sprintf("%s12", url)
  }else{
    fmt.Println("Sun set - moving camera to garage")
    url = fmt.Sprintf("%s17", url)
  }

  r := restapi.NewGet("sunriseset", url)

  r.SetBasicAccessToken(getToken(false))

  //r.DebugOn()

  restapi.TurnOffCertValidation()

  if(r.Send()){

    //r.Dump()

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
  fmt.Println("day - Move camera to day possition\n");
  fmt.Println("night - Move camera to night possition\n");
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

  logmsg.SetLogLevel(logmsg.Debug03)

  logmsg.SetLogFile(logfile)

  args := os.Args;

  if(len(args) >=2){ // being used for other reasons than moving the camera

    switch args[1]{
      
      case "auth":
        buildAuthToken(args[2])

      case "show":
        fmt.Println("Token: ", getToken(false))

      case "showdecoded":
        fmt.Println("Token: ", getToken(true))


      case "day":
        fmt.Println("day move")
        moveCamera(true)

      case "night":
        fmt.Println("night move")
        moveCamera(false)

      default:
        help()

    }

    os.Exit(0)

  }

  if(getToken(false) == "notset"){
    fmt.Println("Need auth token saved")
    os.Exit(0)
  }

  now := time.Now()

  timeRise, timeSet := getSunTimes()

  logmsg.Print(logmsg.Info, "Sunrise: ", timeRise)
  logmsg.Print(logmsg.Info, "Sunset: ", timeSet)

  // test for after sunrise
  if(now.After(timeRise) && now.Before(timeSet)){

    logmsg.Print(logmsg.Info, "Sun has rose")

    moveCamera(true)
  }else{

    moveCamera(false)
    logmsg.Print(logmsg.Info, "Sun has set")

  }



}
