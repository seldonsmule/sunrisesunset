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
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/restapi"

)

type RiseSet struct {

  timeRise time.Time
  timeSet time.Time

}

func getSunTimes() (time.Time, time.Time) {

  var astroMap map[string]interface{}

  info, statErr := os.Stat("/tmp/sun.json")

  if(statErr == nil){

    fmt.Println("Last sun.json write - ",info.ModTime())

    today := time.Now()

    if(today.Day() != info.ModTime().Day()){
      fmt.Println("not today - delete sun.json")
      os.Remove("/tmp/sun.json")
    }

  }

  jsonReadFile, openErr := os.Open("/tmp/sun.json")

  if(openErr == nil){

    //fmt.Println("found the file");

    byteValue, _ := ioutil.ReadAll(jsonReadFile)

    json.Unmarshal([]byte(byteValue), &astroMap)

    jsonReadFile.Close()

    //fmt.Println(astroMap["sunset"])

  }else{

  r := restapi.NewGet("sunriseset", "https://weather.cit.api.here.com/weather/1.0/report.json?product=forecast_astronomy&name=DC&app_id=DemoAppId01082013GAL&app_code=AJKnXv84fjrb0KIHawS0Tg")


//  r.DebugOn()

  r.HasInnerMap("astronomy")

  if(r.Send()){

 //   r.Dump()

  }

/*
  fmt.Printf("--------------------------\n")
  fmt.Printf("astronomy[%s]\n", r.GetValue("astronomy"))
  fmt.Printf("--------------------------\n")
*/

// here is the deal - this stupid thing is an array of maps, not
// something considered in my original Telsa use case.  here is how to
// get the data...

// 1. Get the info as an array

  astroArray := restapi.CastArray(r.GetValue("astronomy"))

/*
  fmt.Printf("ArrayLength[%d]\n", len(astroArray))
*/

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

  jsonWriteFile, err := os.Create("/tmp/sun.json")

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

func main() {

  logmsg.SetLogFile("sun.log")

  now := time.Now()

  timeRise, timeSet := getSunTimes()


/*
  fmt.Println("Now:", now)
  fmt.Println("Now:", now.Unix())
  fmt.Println("Sunrise:", timeRise)
  fmt.Println("Sunrise:", timeRise.Unix())
  fmt.Println("Sunset:", timeSet)
  fmt.Println("Sunset:", timeSet.Unix())
*/
 
  fmt.Println("Sunrise:", timeRise)
  fmt.Println("Sunset:", timeSet)

  // test for after sunrise
  if(now.After(timeRise) && now.Before(timeSet)){

    fmt.Println("Sun has rose")

  }else{

    fmt.Println("Sun has set")

  }



}
