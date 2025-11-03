package config

import "testing"

func TestDotNotationMap(t *testing.T) {
	t.Run("get and set simple values", func(t *testing.T) {
		m := newDotNotationMap()
		m.Set("key", "value")
		val := m.Get("key")
		if val != "value" {
			t.Errorf("Get() = %v, want %v", val, "value")
		}
	})

	t.Run("get and set nested values", func(t *testing.T) {
		m := newDotNotationMap()
		m.Set("mqtt.broker", "tcp://localhost:1883")
		m.Set("mqtt.qos", 1)
		m.Set("device.id", "test-device")

		if m.Get("mqtt.broker") != "tcp://localhost:1883" {
			t.Errorf("Get(mqtt.broker) = %v, want %v", m.Get("mqtt.broker"), "tcp://localhost:1883")
		}
		if m.Get("mqtt.qos") != 1 {
			t.Errorf("Get(mqtt.qos) = %v, want %v", m.Get("mqtt.qos"), 1)
		}
		if m.Get("device.id") != "test-device" {
			t.Errorf("Get(device.id) = %v, want %v", m.Get("device.id"), "test-device")
		}
	})

	t.Run("get non-existent key", func(t *testing.T) {
		m := newDotNotationMap()
		val := m.Get("non.existent.key")
		if val != nil {
			t.Errorf("Get(non.existent.key) = %v, want nil", val)
		}
	})

	t.Run("deeply nested values", func(t *testing.T) {
		m := newDotNotationMap()
		m.Set("a.b.c.d", "deep value")
		val := m.Get("a.b.c.d")
		if val != "deep value" {
			t.Errorf("Get(a.b.c.d) = %v, want %v", val, "deep value")
		}
	})

	t.Run("override existing value", func(t *testing.T) {
		m := newDotNotationMap()
		m.Set("key", "old value")
		m.Set("key", "new value")
		val := m.Get("key")
		if val != "new value" {
			t.Errorf("Get(key) = %v, want %v", val, "new value")
		}
	})

	t.Run("get all as map", func(t *testing.T) {
		m := newDotNotationMap()
		m.Set("a.b", "value1")
		m.Set("a.c", "value2")
		m.Set("d", "value3")

		allMap := m.GetAllAsMap()
		if allMap == nil {
			t.Error("GetAllAsMap() returned nil")
		}

		aMap, ok := allMap["a"].(map[string]interface{})
		if !ok {
			t.Error("GetAllAsMap() 'a' is not a map")
		}
		if aMap["b"] != "value1" {
			t.Errorf("GetAllAsMap() a.b = %v, want %v", aMap["b"], "value1")
		}
		if aMap["c"] != "value2" {
			t.Errorf("GetAllAsMap() a.c = %v, want %v", aMap["c"], "value2")
		}
		if allMap["d"] != "value3" {
			t.Errorf("GetAllAsMap() d = %v, want %v", allMap["d"], "value3")
		}
	})

	t.Run("nested map structure", func(t *testing.T) {
		m := newDotNotationMap()
		m.Set("log.level", "info")
		m.Set("log.source.enabled", true)
		m.Set("log.source.relative", false)

		if m.Get("log.level") != "info" {
			t.Errorf("Get(log.level) = %v, want %v", m.Get("log.level"), "info")
		}
		if m.Get("log.source.enabled") != true {
			t.Errorf("Get(log.source.enabled) = %v, want %v", m.Get("log.source.enabled"), true)
		}
		if m.Get("log.source.relative") != false {
			t.Errorf("Get(log.source.relative) = %v, want %v", m.Get("log.source.relative"), false)
		}
	})
}
