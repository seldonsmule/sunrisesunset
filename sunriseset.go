package main

/* simple get sunrise/set info */

import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/restapi"

)

type RiseSet struct {

  timeRise time.Time
  timeSet time.Time

}

func getSunTimes() (time.Time, time.Time) {

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

  astroArray := r.CastArray(r.GetValue("astronomy"))

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
  astroMap := r.CastMap(astroArray[0])

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

  sunrise :=  r.CastString(astroMap["sunrise"])
  sunset  :=  r.CastString(astroMap["sunset"])

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

  var myTimes RiseSet

  logmsg.SetLogFile("sun.log")

  now := time.Now()

  //timeRise, timeSet := getSunTimes()
  myTimes.timeRise, myTimes.timeSet = getSunTimes()


/*
  fmt.Println("Now:", now)
  fmt.Println("Now:", now.Unix())
  fmt.Println("Sunrise:", timeRise)
  fmt.Println("Sunrise:", timeRise.Unix())
  fmt.Println("Sunset:", timeSet)
  fmt.Println("Sunset:", timeSet.Unix())
*/
 
  // test for after sunrise
  if(now.After(myTimes.timeRise) && now.Before(myTimes.timeSet)){

    fmt.Println("Sun has rose")

  }else{

    fmt.Println("Sun has set")

  }



}
