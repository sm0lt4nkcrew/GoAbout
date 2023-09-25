// GoAbout - track information about IP using ip-api.com API.
// by ...
package main

import (
  "net/http"
  "net/url"
  "log"
  "fmt"
  "encoding/json"
  "os"
  "time"
)

// API response structore.
// In order to change them, change filed number in compile_url function.
// Number can be generated via API.
type API_response struct {
  Status    string
  Country   string
  City      string
  Lat       float64
  Lon       float64
  Isp       string
  Org       string
  As        string
  Mobile    bool
  Proxy     bool
  Hosting   bool
  Query     string
}

// Read args & orchestrate work. 2 arguments are required: mode ip
// mode can be set to: 
//    "show" to print information only once
//    "monitor" to print information constantly every 5 sec.
// ip should be either a valid ip address or keyword "me" for current connetion.
func main() {
  if os.Args[1] == "show" || os.Args[1] == "monitor" {
    full_url := compile_url(os.Args[2])
      for {
        data, err := make_request(full_url)
        if err != nil {
          log.Fatalln(err)
        }
        if os.Args[1] == "monitor" {
          fmt.Println("\033[H\033[2j")
        }

        print_data(data)

        if os.Args[1] == "show" {
         break
        } else {
          time.Sleep(5 * time.Second)
        }
      }
  } else {
    fmt.Printf("Error.")
  }
}

// Parse url with ip as endpoint.
func compile_url(endpoint string) string {
  api_url := "http://ip-api.com"
  base_url, base_err := url.Parse(api_url)
  if base_err != nil {
    log.Fatalln(base_err)
  }

  if endpoint != "me" {
    base_url.Path += endpoint + "/"
  }
  base_url.Path += "json"

  params, params_err := url.ParseQuery("fields=17002193")
  if params_err != nil {
    log.Fatalln(params_err)
  }
  base_url.RawQuery = params.Encode()
  return base_url.String()
}

// Make API request.
func make_request(url string) (*API_response, error) {
  res, resp_err := http.Get(url)
  if resp_err != nil {
    log.Fatalln(resp_err)
  }

  var data API_response
  read_err := json.NewDecoder(res.Body).Decode(&data)
  if read_err != nil {
    log.Fatalln(read_err)
  }

  defer res.Body.Close()
  
  return &data, nil
}

// Print data.
func print_data(resp *API_response) {
  t := time.Now()
  fmt.Printf("Check time: %d:%d:%d\n", t.Hour(), t.Minute(), t.Second())
  fmt.Println("IP:", resp.Query)
  fmt.Printf("Location: %s, %s (%f, %f)\n", resp.Country, resp.City, resp.Lat, resp.Lon)
  fmt.Printf("ISP: %s for %s as %s\n", resp.Isp, resp.Org, resp.As)
  fmt.Println("Mobile:", resp.Mobile)
  fmt.Println("Proxy:", resp.Proxy)
  fmt.Println("Hosting:", resp.Hosting)
}


