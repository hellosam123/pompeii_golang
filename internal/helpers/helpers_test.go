package helpers

import (
	"testing"
)

func TestGetAllVocab(t *testing.T) {
	vocabs, err := GetAllVocab()
	if err != nil {
		t.Fatalf("GetAllVocab failed: %v", err)
	}

	t.Log(vocabs)
}

func TestGetVocabByGroup(t *testing.T) {
	vocabs, err := GetVocabByGroup("Latin IGCSE")
	if err != nil {
		t.Fatalf("GetVocabByGroups failed: %v", err)
	}

	t.Log(vocabs)
}

func TestGetVocabByID(t *testing.T) {
	vocab, err := GetVocabByID(4)
	if err != nil {
		t.Fatalf("GetVocabByID failed: %v", err)
	}

	t.Log(vocab)
}

func TestGetVocabByIDAllTranslations(t *testing.T) {
	vocab, err := GetVocabByIDAllTranslations(4)
	if err != nil {
		t.Fatalf("GetVocabByIDAllTranslations failed: %v", err)
	}

	t.Log(vocab)
}

func TestGetRandomVocabByGroup(t *testing.T) {
	vocab, err := GetRandomVocabByGroup("Latin IGCSE")
	if err != nil {
		t.Fatalf("GetVocabByGroups failed: %v", err)
	}

	t.Log(vocab)
	t.Errorf("clear cache")
}

func TestCheckAnswer(t *testing.T) {
	correct, err := CheckAnswer(4, "I am accepting")
	if err != nil {
		t.Fatalf("GetVocabByGroups failed: %v", err)
	}

	t.Log(correct)
}

func BenchmarkGetRandomVocabByGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetRandomVocabByGroup("Latin IGCSE")
		if err != nil {
			b.Fatalf("GetRandomVocabByGroup failed: %v", err)
		}
	}
}
