package main

/* simple get sunrise/set info */

import (
        "fmt"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/restapi"

)

func main() {

  logmsg.SetLogFile("sun.log")

  fmt.Println("Sunriseset")


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
  astroMap := r.CastMap(astroArray[1])

// used this to figure out the names to extract the sunset/rise info - unommit to see
/*
  for k, v := range astroMap {
    fmt.Println(k, "=", v)
  } // end for loop
*/


// 3. get teh value in the map
  fmt.Printf("--------------------------\n")
  fmt.Printf("sunset[%s]\n", astroMap["sunset"])
  fmt.Printf("sunrise[%s]\n", astroMap["sunrise"])



 
}
