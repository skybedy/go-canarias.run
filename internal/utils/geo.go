package utils

import "strings"

type IslandRule struct {
	Name     string
	Keywords []string
}

var islandRules = []IslandRule{
	{
		Name: "La Palma",
		Keywords: []string{
			"la palma", "transvulcania", "los llanos", "tazacorte", "fuencaliente",
			"breña alta", "breña baja", "puntagorda", "garafía", "tijarafe",
			"barlovento", "puntallana", "el paso", "taburiente", "reventon", "roque de los muchachos",
			"santa cruz de la palma", "puntallana", "pazur", "teneguia", "mágica",
		},
	},
	{
		Name: "Tenerife",
		Keywords: []string{
			"tenerife", "santa cruz", "la laguna", "arona", "adeje", "puerto de la cruz",
			"realejos", "granadilla", "orotava", "icod", "candelaria", "tacoronte",
			"guía de isora", "güímar", "guimar", "el rosario", "san miguel", "sauzal", "tegueste",
			"victoria", "matanza", "arico", "santiago del teide", "vilaflor", "fasnia",
			"buenavista", "los silos", "garachico", "tanque", "teide", "anaga", "mercedes",
			"bluetrail", "tfe", "el medano", "pinolere", "tinguaro", "carboneras", "esperanza",
			"tamaimo", "la cruz santa", "los abrigos", "ravelo", "asomadero", "vallivana",
			"maretas", "brutaltrail", "don leandro", "anaga", "viera", "clavería",
		},
	},
	{
		Name: "Gran Canaria",
		Keywords: []string{
			"gran canaria", "las palmas", "telde", "santa lucía", "santa lucia", "san bartolomé de tirajana",
			"arucas", "agüimes", "aguimes", "ingenio", "gáldar", "galdar", "mogán", "mogan", "santa brígida",
			"teror", "valsequillo", "aldea de san nicolás", "moya", "firgas", "agaete", "san mateo",
			"artenara", "tejeda", "valleseco", "roque nublo", "tamadaba", "maspalomas", "vecindario",
			"transgrancanaria", "binter night run", "animas", "flick", "freedon", "gc", "acebuches",
			"aguas de teror", "guayadeque", "arintegui", "pino trail", "aguas de teror", "confital",
			"lpa trail", "entre cortijos", "salesianas", "stay alive", "360 the challenge", "pinar de tamadaba",
		},
	},
	{
		Name: "Lanzarote",
		Keywords: []string{
			"lanzarote", "arrecife", "teguise", "tías", "tias", "san bartolomé", "yaiza", "tinajo",
			"haría", "haria", "famara", "timanfaya", "la geria", "costa teguise", "playa blanca", "puerto del carmen",
			"ironman lanzarote", "clm", "lzt", "pinto", "isla graciosa", "malpaso", "tinajo you trail",
			"wine run", "caldera de el cuchillo", "montaña blanca", "teseguite", "san rafael", "haria titan",
			"mascaritas", "valle de malpaso", "teguise", "haria",
		},
	},
	{
		Name: "Fuerteventura",
		Keywords: []string{
			"fuerteventura", "puerto del rosario", "la oliva", "pájara", "pajara", "tuineje",
			"antigua", "betancuria", "corralejo", "jandía", "jandia", "el cotillo", "morro jable",
			"fv", "majorera", "carnavalera", "molinos", "la lajita", "betancuria", "pajara", "antigua",
			"caleta de fuste", "costa calma", "corralejo", "grandes playas", "isla de lobos",
		},
	},
	{
		Name: "La Gomera",
		Keywords: []string{
			"la gomera", "san sebastián de la gomera", "san sebastian de la gomera", "valle gran rey", "alajeró", "alajero", "hermigua",
			"agulo", "vallehermoso", "garajonay", "gomera", "hermigua", "ipalán", "ipalan", "gomera paradise",
			"majalca", "cedro", "san sebastian",
		},
	},
	{
		Name: "El Hierro",
		Keywords: []string{
			"el hierro", "valverde", "frontera", "el pinar", "restinga", "meridiano", "mar de las calmas", "hierro",
		},
	},
	{
		Name: "La Graciosa",
		Keywords: []string{
			"la graciosa", "caleta de sebo",
		},
	},
}

func IdentifyIsland(text ...string) string {
	combined := strings.Join(text, " ")
	fullText := strings.ToLower(combined)

	for _, rule := range islandRules {
		for _, kw := range rule.Keywords {
			if strings.Contains(fullText, kw) {
				return rule.Name
			}
		}
	}

	return "Canarias"
}
