package fxtx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Gpx
type Gpx struct {
	XMLName xml.Name `xml:"gpx"`
	Wpts    []Wpt    `xml:"wpt"`
}

// Wpt
type Wpt struct {
	XMLName xml.Name `xml:"wpt"`
	Lat     float64  `xml:"lat,attr"`
	Lon     float64  `xml:"lon,attr"`
}

// GenCfg
type GenCfg struct {
	Generators []Generator `yaml:"generators"`
}

// Generator
type Generator struct {
	Description      string `yaml:"description"`
	Frequency        int    `yaml:"frequency"`
	WaypointFile     string `yaml:"waypointFile"`
	Template         string `yaml:"template"`
	gpxData          *Gpx
	compiledTemplate *template.Template
}

// GenCfgFromFile creates a configuration file from YAML
func GenCfgFromFile(file string) (*GenCfg, error) {
	ymlData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &GenCfg{}

	err = yaml.Unmarshal([]byte(ymlData), &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Cfg
type Cfg struct {
	GenCfg      *GenCfg
	Destination string
	Timeout     time.Duration
	Logger      *zap.Logger
}

// Fxtx
type Fxtx struct {
	*Cfg
}

// NewFxtx
func NewFxtx(cfg *Cfg) (*Fxtx, error) {

	for i, gen := range cfg.GenCfg.Generators {
		cfg.Logger.Info("Loading Generator", zap.String("description", gen.Description))

		// load waypoints
		cfg.Logger.Info("Loading waypoint file go...", zap.String("waypont_file", gen.WaypointFile))
		gxpData, err := ioutil.ReadFile(gen.WaypointFile)
		if err != nil {
			return nil, fmt.Errorf("error opening waypoint file: %s", err.Error())
		}

		gpx := &Gpx{}
		err = xml.Unmarshal(gxpData, gpx)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling waypoints: %s", err.Error())
		}

		cfg.GenCfg.Generators[i].gpxData = gpx

		// compile template
		cfg.Logger.Info("Compiling template", zap.String("teamplte", gen.Template))
		tmpl, err := template.New("msg_template").Funcs(sprig.TxtFuncMap()).Parse(gen.Template)
		if err != nil {
			return nil, fmt.Errorf("error processing template: %s", err.Error())
		}

		cfg.GenCfg.Generators[i].compiledTemplate = tmpl
	}

	return &Fxtx{cfg}, nil
}

// generate
func (fx *Fxtx) generate(gen Generator, wg *sync.WaitGroup) {
	fx.Logger.Info("Starting Generator", zap.String("generator", gen.Description))

	var max = len(gen.gpxData.Wpts)
	i := 0
	count := 0

	d := net.Dialer{Timeout: fx.Timeout}

	// loop though way-points and generate a message
	for {
		wp := gen.gpxData.Wpts[i]

		// populate parameter map
		params := make(map[string]interface{})
		params["lat"] = wp.Lat
		params["lon"] = wp.Lon

		fx.Logger.Info("Generating message",
			zap.Int("index", i),
			zap.Int("count", count),
		)

		msg := bytes.Buffer{}

		err := gen.compiledTemplate.Execute(&msg, params)
		if err != nil {
			fx.Logger.Error("Error executing message template. Exiting generator.", zap.Error(err))
			break
		}

		msgBytes := msg.Bytes()

		// rendered message
		fx.Logger.Info("Sending rendered message.",
			zap.String("generator", gen.Description),
			zap.ByteString("msg", msgBytes),
		)

		conn, err := d.Dial("tcp", fx.Destination)
		if err != nil {
			fx.Logger.Error("unable to connect",
				zap.String("destination", fx.Destination),
				zap.Error(err))
		}

		_, err = fmt.Fprintf(conn, msg.String()+"\n")
		if err != nil {
			fx.Logger.Error("unable to write ",
				zap.String("destination", fx.Destination),
				zap.Error(err))
		}

		err = conn.Close()
		if err != nil {
			fx.Logger.Error("unable close tcp connect", zap.Error(err))
		}

		// then wait on frequency
		fx.Logger.Info("Wait on interval after send.",
			zap.String("generator", gen.Description),
			zap.Int("frequency", gen.Frequency),
		)

		time.Sleep(time.Duration(gen.Frequency) * time.Second)
		i += 1
		if i >= max {
			i = 0
		}
	}

	wg.Done()
}

// Run
func (fx *Fxtx) Run() {
	var wg sync.WaitGroup

	for _, gen := range fx.GenCfg.Generators {
		wg.Add(1)
		go fx.generate(gen, &wg)
	}

	wg.Wait()
	fx.Logger.Warn("All generators returned.")
}
