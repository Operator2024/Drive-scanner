package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var name = "Drive-scanner"

var (
	version          string
	date             string
	compileDate      time.Time
	goCompileVersion string
)

// Device contains the name of the drive that processed now.
type Device struct {
	Name string `json:"name"`
}

// DeviceInfoRAW contains raw data from smartctl in JSON format.
type DeviceInfoRAW struct {
	Vendor       string     `json:"vendor"`
	Product      string     `json:"product"`
	Model        string     `json:"model_name"`
	SN           string     `json:"serial_number"`
	Rev          string     `json:"revision"`
	FW           string     `json:"firmware_version"`
	Type         int        `json:"rotation_rate"`
	UserCapacity DeviceSize `json:"user_capacity"`
	Size         int        `json:"nvme_total_capacity"`
}

// DeviceSize contains fields related to disk size.
type DeviceSize struct {
	Blocks int `json:"blocks"`
	Bytes  int `json:"bytes"`
}

// DeviceInfo contains already processed disk data
type DeviceInfo struct {
	Vendor  string  `json:"vendor"`
	Product *string `json:"product"`
	Model   string  `json:"model_name"`
	SN      string  `json:"serial_number"`
	Rev     *string `json:"revision"`
	FW      string  `json:"firmware_version"`
	Type    string  `json:"type"`
	Size    int     `json:"size"`
}

// GetDeviceName only returns the drive name.
func GetDeviceName() []byte {
	cmd := exec.Command("smartctl", "--scan", "-j")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}
	return stdout
}

// Usage overrides method 'Usage' from flag.
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)

	flag.PrintDefaults()
}

func main() {
	var buildVersion bool
	flag.BoolVar(&buildVersion, "V", false, "This key allows you to get the current version")

	if version != "" {
		goCompileVersion = runtime.Version()
		compileDate, _ = time.Parse("2006-01-02 03:04:05PM MST", date)
	}
	flag.Usage = Usage
	flag.Parse()

	if buildVersion {
		fmt.Printf("Version: %s, build info: %s [%s]\n", version, goCompileVersion, compileDate.Format("2006-01-02 03:04:05PM MST"))
	} else {
		r, _ := regexp.Compile("\"name\":.{1,},")

		storage := make(map[string][]map[string]DeviceInfo)
		storage[name] = make([]map[string]DeviceInfo, 1)

		var drive Device
		var totalDeviceList = make(map[string]DeviceInfo)
		deviceName := GetDeviceName()

		for _, val := range r.FindAllString(string(deviceName), -1) {
			err := json.Unmarshal([]byte(`{`+val[0:len(val)-1]+`}`), &drive)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			cmd := exec.Command("smartctl", "-i", drive.Name, "-j")
			stdout, err := cmd.Output()

			var d DeviceInfoRAW
			var o DeviceInfo
			err = json.Unmarshal([]byte(stdout), &d)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if d.Model != "" {
				driveVendor := []string{"Apacer", "Samsung", "Seagate", "Toshiba",
					"Apple", "Micron", "Crucial", "Lexar",
					"Kingston", "OCZ", "SanDisk", "WD",
					"Hitachi", "HGST", "PNY", "Corsair", "Intel"}
				for _, v := range driveVendor {
					if strings.Contains(strings.ToLower(d.Model), strings.ToLower(v)) {
						if d.Vendor == "" {
							d.Vendor = v
							d.Model = strings.Replace(strings.ToLower(d.Model), strings.ToLower(v), "", -1)
						}
					}
					if strings.Contains(strings.ToLower(d.Model), strings.ToLower("ST")) {
						if d.Vendor == "" {
							d.Vendor = "Seagate"
						}
					}
				}
			}
			o.Model = strings.TrimLeft(strings.Replace(d.Model, "_", " ", -1), " ")
			o.Vendor = d.Vendor
			o.SN = d.SN
			o.FW = d.FW
			if d.Type == 0 {
				o.Type = "Solid State Drive"
			} else {
				o.Type = "Hard Disk Drive"
			}

			if (d.UserCapacity.Bytes == 0) && (d.Size != 0) {
				// nvme_total_capacity
				o.Size = d.Size
			} else {
				o.Size = d.UserCapacity.Bytes
			}

			if d.Rev != "" {
				o.Rev = &d.Rev
			} else {
				o.Rev = nil
			}

			if d.Product != "" {
				o.Product = &d.Product
			} else {
				o.Product = nil
			}

			totalDeviceList[drive.Name] = o
		}
		storage[name][0] = totalDeviceList
		result, _ := json.Marshal(storage)
		fmt.Printf("%v\n", string(result))
	}
}
