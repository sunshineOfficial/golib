package language

import (
	"testing"
)

func TestTranslitRuToEn(t *testing.T) {
	res := TranslitRuToEn("ЁАБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмнопрстуфхцчшщъыьэюяё")
	if res != "EABVGDEZhZIIKLMNOPRSTUFKhTsChShShchYEIuIaabvgdezhziiklmnoprstufkhtschshshchyeiuiae" {
		t.Error("Error transliterate Ru to En")
	}
}

func TestTranslitEnToRu(t *testing.T) {
	res := TranslitEnToRu("YayaJujuSchschChchShshZhzhYoyoKhkhXxAaBbVvGgDdZzIiJjKkLlMmNnOoPpRrSsTtUuFfCcYyEe")
	if res != "ЯяЮюЩщЧчШшЖжЁёХхКсксАаБбВвГгДдЗзИиЙйКкЛлМмНнОоПпРрСсТтУуФфЦцЫыЕе" {
		t.Error("Error transliterate En to Ru")
	}
}

func TestTranslitEnToRu2(t *testing.T) {
	res := TranslitEnToRu("JajaYuyuJojoHh")
	if res != "ЯяЮюЁёХх" {
		t.Error("Error transliterate En to Ru 2")
	}
}
