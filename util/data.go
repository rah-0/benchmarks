package util

import (
	"math/rand"

	"github.com/google/uuid"
)

type Person struct {
	Name    string
	Age     int
	Surname string
}

func GenerateRandomPersons(count int) []Person {
	persons := make([]Person, count)
	for i := 0; i < count; i++ {
		persons[i] = Person{
			Name:    uuid.NewString(),
			Age:     rand.Intn(100),
			Surname: uuid.NewString(),
		}
	}
	return persons
}

type Unreal struct {
	Prop1   string
	Prop2   string
	Prop3   string
	Prop4   string
	Prop5   string
	Prop6   string
	Prop7   string
	Prop8   string
	Prop9   string
	Prop10  string
	Prop11  string
	Prop12  string
	Prop13  string
	Prop14  string
	Prop15  string
	Prop16  string
	Prop17  string
	Prop18  string
	Prop19  string
	Prop20  string
	Prop21  string
	Prop22  string
	Prop23  string
	Prop24  string
	Prop25  string
	Prop26  string
	Prop27  string
	Prop28  string
	Prop29  string
	Prop30  string
	Prop31  string
	Prop32  string
	Prop33  string
	Prop34  string
	Prop35  string
	Prop36  string
	Prop37  string
	Prop38  string
	Prop39  string
	Prop40  string
	Prop41  string
	Prop42  string
	Prop43  string
	Prop44  string
	Prop45  string
	Prop46  string
	Prop47  string
	Prop48  string
	Prop49  string
	Prop50  string
	Prop51  string
	Prop52  string
	Prop53  string
	Prop54  string
	Prop55  string
	Prop56  string
	Prop57  string
	Prop58  string
	Prop59  string
	Prop60  string
	Prop61  string
	Prop62  string
	Prop63  string
	Prop64  string
	Prop65  string
	Prop66  string
	Prop67  string
	Prop68  string
	Prop69  string
	Prop70  string
	Prop71  string
	Prop72  string
	Prop73  string
	Prop74  string
	Prop75  string
	Prop76  string
	Prop77  string
	Prop78  string
	Prop79  string
	Prop80  string
	Prop81  string
	Prop82  string
	Prop83  string
	Prop84  string
	Prop85  string
	Prop86  string
	Prop87  string
	Prop88  string
	Prop89  string
	Prop90  string
	Prop91  string
	Prop92  string
	Prop93  string
	Prop94  string
	Prop95  string
	Prop96  string
	Prop97  string
	Prop98  string
	Prop99  string
	Prop100 string
}

func GenerateRandomUnreals(count int) []Unreal {
	unreals := make([]Unreal, count)
	for i := 0; i < count; i++ {
		unreals[i] = Unreal{
			Prop1:   uuid.NewString(),
			Prop2:   uuid.NewString(),
			Prop3:   uuid.NewString(),
			Prop4:   uuid.NewString(),
			Prop5:   uuid.NewString(),
			Prop6:   uuid.NewString(),
			Prop7:   uuid.NewString(),
			Prop8:   uuid.NewString(),
			Prop9:   uuid.NewString(),
			Prop10:  uuid.NewString(),
			Prop11:  uuid.NewString(),
			Prop12:  uuid.NewString(),
			Prop13:  uuid.NewString(),
			Prop14:  uuid.NewString(),
			Prop15:  uuid.NewString(),
			Prop16:  uuid.NewString(),
			Prop17:  uuid.NewString(),
			Prop18:  uuid.NewString(),
			Prop19:  uuid.NewString(),
			Prop20:  uuid.NewString(),
			Prop21:  uuid.NewString(),
			Prop22:  uuid.NewString(),
			Prop23:  uuid.NewString(),
			Prop24:  uuid.NewString(),
			Prop25:  uuid.NewString(),
			Prop26:  uuid.NewString(),
			Prop27:  uuid.NewString(),
			Prop28:  uuid.NewString(),
			Prop29:  uuid.NewString(),
			Prop30:  uuid.NewString(),
			Prop31:  uuid.NewString(),
			Prop32:  uuid.NewString(),
			Prop33:  uuid.NewString(),
			Prop34:  uuid.NewString(),
			Prop35:  uuid.NewString(),
			Prop36:  uuid.NewString(),
			Prop37:  uuid.NewString(),
			Prop38:  uuid.NewString(),
			Prop39:  uuid.NewString(),
			Prop40:  uuid.NewString(),
			Prop41:  uuid.NewString(),
			Prop42:  uuid.NewString(),
			Prop43:  uuid.NewString(),
			Prop44:  uuid.NewString(),
			Prop45:  uuid.NewString(),
			Prop46:  uuid.NewString(),
			Prop47:  uuid.NewString(),
			Prop48:  uuid.NewString(),
			Prop49:  uuid.NewString(),
			Prop50:  uuid.NewString(),
			Prop51:  uuid.NewString(),
			Prop52:  uuid.NewString(),
			Prop53:  uuid.NewString(),
			Prop54:  uuid.NewString(),
			Prop55:  uuid.NewString(),
			Prop56:  uuid.NewString(),
			Prop57:  uuid.NewString(),
			Prop58:  uuid.NewString(),
			Prop59:  uuid.NewString(),
			Prop60:  uuid.NewString(),
			Prop61:  uuid.NewString(),
			Prop62:  uuid.NewString(),
			Prop63:  uuid.NewString(),
			Prop64:  uuid.NewString(),
			Prop65:  uuid.NewString(),
			Prop66:  uuid.NewString(),
			Prop67:  uuid.NewString(),
			Prop68:  uuid.NewString(),
			Prop69:  uuid.NewString(),
			Prop70:  uuid.NewString(),
			Prop71:  uuid.NewString(),
			Prop72:  uuid.NewString(),
			Prop73:  uuid.NewString(),
			Prop74:  uuid.NewString(),
			Prop75:  uuid.NewString(),
			Prop76:  uuid.NewString(),
			Prop77:  uuid.NewString(),
			Prop78:  uuid.NewString(),
			Prop79:  uuid.NewString(),
			Prop80:  uuid.NewString(),
			Prop81:  uuid.NewString(),
			Prop82:  uuid.NewString(),
			Prop83:  uuid.NewString(),
			Prop84:  uuid.NewString(),
			Prop85:  uuid.NewString(),
			Prop86:  uuid.NewString(),
			Prop87:  uuid.NewString(),
			Prop88:  uuid.NewString(),
			Prop89:  uuid.NewString(),
			Prop90:  uuid.NewString(),
			Prop91:  uuid.NewString(),
			Prop92:  uuid.NewString(),
			Prop93:  uuid.NewString(),
			Prop94:  uuid.NewString(),
			Prop95:  uuid.NewString(),
			Prop96:  uuid.NewString(),
			Prop97:  uuid.NewString(),
			Prop98:  uuid.NewString(),
			Prop99:  uuid.NewString(),
			Prop100: uuid.NewString(),
		}
	}
	return unreals
}
