package msg

import (
	"testing"
)

func messageBuffer() []byte {
	return make([]byte, binarySize)
}

func TestMarshallDoesNotDestroyMesssage(t *testing.T) {
	ref := &Message{Kind: 1, Price: 2, Amount: 3, StockId: 4, TraderId: 5, TradeId: 6}
	m1 := &Message{}
	*m1 = *ref
	b := messageBuffer()
	if err := Marshal(b, m1); err != nil {
		t.Errorf("Unexpected marshalling error %s", err.Error())
	}
	if *m1 != *ref {
		t.Errorf("Expected to find %v, found %v instead. Marshalled from %v", ref, m1, b)
	}
}

func TestMarshallUnMarshalPairsProducesSameMessage(t *testing.T) {
	m1 := &Message{Kind: 1, Price: 2, Amount: 3, StockId: 4, TraderId: 5, TradeId: 6}
	b := messageBuffer()
	if err := Marshal(b, m1); err != nil {
		t.Errorf("Unexpected marshalling error %s", err.Error())
	}
	m2 := &Message{}
	if err := Unmarshal(b, m2); err != nil {
		t.Errorf("Unexpected unmarshalling error %s", err.Error())
	}
	if *m2 != *m1 {
		t.Errorf("Expected to find %v, found %v instead. Marshalled from %v", m1, m2, b)
	}
}

func TestMarshalWithSmallBufferErrors(t *testing.T) {
	m1 := &Message{Kind: 1, Price: 2, Amount: 3, StockId: 4, TraderId: 5, TradeId: 6}
	b := make([]byte, binarySize-1)
	if err := Marshal(b, m1); err == nil {
		t.Error("Expected marshalling error. Found none")
	}
}

func TestMarshalWithLargeBufferErrors(t *testing.T) {
	m1 := &Message{Kind: 1, Price: 2, Amount: 3, StockId: 4, TraderId: 5, TradeId: 6}
	b := make([]byte, binarySize+1)
	if err := Marshal(b, m1); err == nil {
		t.Error("Expected marshalling error. Found none")
	}
}

func TestUnmarshalWithSmallBufferErrors(t *testing.T) {
	m1 := &Message{}
	b := make([]byte, binarySize-1)
	if err := Unmarshal(b, m1); err == nil {
		t.Error("Expected marshalling error. Found none")
	}
}

func TestUnmarshalWithLargeBufferErrors(t *testing.T) {
	m1 := &Message{}
	b := make([]byte, binarySize+1)
	if err := Unmarshal(b, m1); err == nil {
		t.Error("Expected marshalling error. Found none")
	}
}
