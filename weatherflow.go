package main

// this was from the weatherflow api

type BetterForcast struct {
	CurrentConditions struct {
		AirDensity                     float64 `json:"air_density"`
		AirTemperature                 float64 `json:"air_temperature"`
		Brightness                     int     `json:"brightness"`
		Conditions                     string  `json:"conditions"`
		DeltaT                         float64 `json:"delta_t"`
		DewPoint                       float64 `json:"dew_point"`
		FeelsLike                      float64 `json:"feels_like"`
		Icon                           string  `json:"icon"`
		IsPrecipLocalDayRainCheck      bool    `json:"is_precip_local_day_rain_check"`
		LightningStrikeCountLast1Hr    int     `json:"lightning_strike_count_last_1hr"`
		LightningStrikeCountLast3Hr    int     `json:"lightning_strike_count_last_3hr"`
		LightningStrikeLastDistance    int     `json:"lightning_strike_last_distance"`
		LightningStrikeLastDistanceMsg string  `json:"lightning_strike_last_distance_msg"`
		LightningStrikeLastEpoch       int     `json:"lightning_strike_last_epoch"`
		PrecipAccumLocalDay            float64 `json:"precip_accum_local_day"`
		PrecipAccumLocalYesterday      int     `json:"precip_accum_local_yesterday"`
		PrecipMinutesLocalDay          int     `json:"precip_minutes_local_day"`
		PrecipMinutesLocalYesterday    int     `json:"precip_minutes_local_yesterday"`
		PressureTrend                  string  `json:"pressure_trend"`
		RelativeHumidity               int     `json:"relative_humidity"`
		SeaLevelPressure               float64 `json:"sea_level_pressure"`
		SolarRadiation                 int     `json:"solar_radiation"`
		StationPressure                float64 `json:"station_pressure"`
		Time                           int     `json:"time"`
		Uv                             int     `json:"uv"`
		WetBulbGlobeTemperature        float64 `json:"wet_bulb_globe_temperature"`
		WetBulbTemperature             float64 `json:"wet_bulb_temperature"`
		WindAvg                        float64 `json:"wind_avg"`
		WindDirection                  int     `json:"wind_direction"`
		WindDirectionCardinal          string  `json:"wind_direction_cardinal"`
		WindGust                       float64 `json:"wind_gust"`
	} `json:"current_conditions"`
	Forecast struct {
		Daily []struct {
			AirTempHigh       float64 `json:"air_temp_high"`
			AirTempLow        float64 `json:"air_temp_low"`
			Conditions        string  `json:"conditions"`
			DayNum            int     `json:"day_num"`
			DayStartLocal     int     `json:"day_start_local"`
			Icon              string  `json:"icon"`
			MonthNum          int     `json:"month_num"`
			PrecipIcon        string  `json:"precip_icon"`
			PrecipProbability int     `json:"precip_probability"`
			PrecipType        string  `json:"precip_type"`
			Sunrise           int     `json:"sunrise"`
			Sunset            int     `json:"sunset"`
		} `json:"daily"`
		Hourly []struct {
			AirTemperature        float64 `json:"air_temperature"`
			Conditions            string  `json:"conditions"`
			FeelsLike             float64 `json:"feels_like"`
			Icon                  string  `json:"icon"`
			LocalDay              int     `json:"local_day"`
			LocalHour             int     `json:"local_hour"`
			Precip                int     `json:"precip"`
			PrecipIcon            string  `json:"precip_icon"`
			PrecipProbability     int     `json:"precip_probability"`
			PrecipType            string  `json:"precip_type"`
			RelativeHumidity      int     `json:"relative_humidity"`
			SeaLevelPressure      float64 `json:"sea_level_pressure"`
			Time                  int     `json:"time"`
			Uv                    float64 `json:"uv"`
			WindAvg               float64 `json:"wind_avg"`
			WindDirection         int     `json:"wind_direction"`
			WindDirectionCardinal string  `json:"wind_direction_cardinal"`
			WindGust              float64 `json:"wind_gust"`
		} `json:"hourly"`
	} `json:"forecast"`
	Latitude           float64 `json:"latitude"`
	LocationName       string  `json:"location_name"`
	Longitude          float64 `json:"longitude"`
	SourceIDConditions int     `json:"source_id_conditions"`
	Status             struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"status"`
	Timezone              string `json:"timezone"`
	TimezoneOffsetMinutes int    `json:"timezone_offset_minutes"`
	Units                 struct {
		UnitsAirDensity     string `json:"units_air_density"`
		UnitsBrightness     string `json:"units_brightness"`
		UnitsDistance       string `json:"units_distance"`
		UnitsOther          string `json:"units_other"`
		UnitsPrecip         string `json:"units_precip"`
		UnitsPressure       string `json:"units_pressure"`
		UnitsSolarRadiation string `json:"units_solar_radiation"`
		UnitsTemp           string `json:"units_temp"`
		UnitsWind           string `json:"units_wind"`
	} `json:"units"`
}
