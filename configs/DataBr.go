package configs

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)




type DataBr time.Time

func (d *DataBr) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
    // Layout brasileiro: dia/mês/ano
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return fmt.Errorf("formato de data inválido. Use DD/MM/YYYY")
	}
	*d = DataBr(t)
	return nil
}

// Opcional: Para devolver o JSON no formato brasileiro também
func (d *DataBr) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*d).Format("02/01/2006"))), nil
}

func NewDataBrPtr(t time.Time) *DataBr {
    d := DataBr(t)
    return &d
}

func (d *DataBr) Time() time.Time {
    return time.Time(*d)
}

func (d *DataBr) IsZero() bool {
    return time.Time(*d).IsZero()
}


func (d DataBr) Value() (driver.Value, error) {
    // Retorna a data nativa do Go (time.Time) para o banco entender
    // Caso seu DataBr guarde a data em um campo interno, use d.SeuCampoDeData
    // Se DataBr for apenas um alias (type DataBr time.Time), faça: return time.Time(d), nil
    
    // Verifique se a data está zerada, se necessário retornar NULL
    if d.IsZero() {
        return nil, nil
    }

	return time.Time(d), nil
}
