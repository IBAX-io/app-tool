package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func readConfig(out *exportFile) {
	config := exportFile{}
	absConfPath, _ := filepath.Abs(inputName)
	bs, err := os.ReadFile(filepath.Join(absConfPath, configName))
	if err != nil {
		if debug {
			fmt.Println("config file not found. used default values")
		}
		return
	}
	_ = json.Unmarshal(bs, &config)
	if len(config.Name) > 0 {
		importNew = true
		out.Name = config.Name
		out.Conditions = config.Conditions
	}
	if len(config.Snippets) > 0 {
		for c := range config.Snippets {
			for o := range out.Snippets {
				if config.Snippets[c].Name == out.Snippets[o].Name {
					out.Snippets[o].Conditions = config.Snippets[c].Conditions
				}
			}
		}
	}
	if len(config.Contracts) > 0 {
		for c := range config.Contracts {
			for o := range out.Contracts {
				if config.Contracts[c].Name == out.Contracts[o].Name {
					out.Contracts[o].Conditions = config.Contracts[c].Conditions
					out.Contracts[o].Confirmation = config.Contracts[c].Confirmation
				}
			}
		}
	}
	if len(config.Menus) > 0 {
		for c := range config.Menus {
			for o := range out.Menus {
				if config.Menus[c].Name == out.Menus[o].Name {
					out.Menus[o].Conditions = config.Menus[c].Conditions
				}
			}
		}
	}
	if len(config.Pages) > 0 {
		for c := range config.Pages {
			for o := range out.Pages {
				if config.Pages[c].Name == out.Pages[o].Name {
					out.Pages[o].Conditions = config.Pages[c].Conditions
					if len(config.Pages[c].Menu) > 0 {
						out.Pages[o].Menu = config.Pages[c].Menu
					}
				}
			}
		}
	}
	if len(config.Tables) > 0 {
		for c := range config.Tables {
			for o := range out.Tables {
				if config.Tables[c].Name == out.Tables[o].Name {
					out.Tables[o].Permissions = config.Tables[c].Permissions
				}
			}
		}
	}
	if len(config.Parameters) > 0 {
		for c := range config.Parameters {
			for o := range out.Parameters {
				if config.Parameters[c].Name == out.Parameters[o].Name {
					out.Parameters[o].Conditions = config.Parameters[c].Conditions
				}
			}
		}
	}
	return
}

func writeConfig(bs []byte) {
	conf := configFile{}
	if importNew {
		tempConf := dataConf{}
		if err := json.Unmarshal(bs, &tempConf); err != nil {
			fmt.Println("unmarshal config file error:", err)
			return
		}
		conf = convertDataConf(tempConf)
	} else {
		if err := json.Unmarshal(bs, &conf); err != nil {
			fmt.Println("unmarshal config file error:", err)
			return
		}
	}
	if bs, err := json.MarshalIndent(conf, "", "    "); err == nil {
		writeFileString(configName, string(bs))
	}
}
func convertDataConf(conf dataConf) (res configFile) {
	res.Name = conf.Name
	res.Conditions = conf.Conditions
	for _, item := range conf.Data {
		switch item.Type {
		case typeSnippet:
			item.Type = ""
			res.Snippets = append(res.Snippets, item)
		case typeMenu:
			item.Type = ""
			res.Menus = append(res.Menus, item)
		case typePage:
			item.Type = ""
			res.Pages = append(res.Pages, item)
		case typeParam:
			item.Type = ""
			res.Params = append(res.Params, item)
		case typeCon:
			item.Type = ""
			res.Contracts = append(res.Contracts, item)
		case typeTable:
			item.Type = ""
			res.Tables = append(res.Tables, item)
		}
	}
	return res
}
