package builtin

import (
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/goland/ast"
)

func init() {
	instanceMethods := NewMethodMap()
	instanceMethods.Set(
		"getErrors",
		[]*Method{
			CreateMethod(
				"getErrors",
				[]string{"List"},
				[]ast.Node{},
				func(this *Object, params []*Object, extra map[string]interface{}) interface{} {
					return this.Extra["errors"]
				},
			),
		},
	)
	instanceMethods.Set(
		"getId",
		[]*Method{
			CreateMethod(
				"getId",
				[]string{"String"},
				[]ast.Node{},
				func(this *Object, params []*Object, extra map[string]interface{}) interface{} {
					return this.Extra["id"]
				},
			),
		},
	)
	instanceMethods.Set(
		"isSuccess",
		[]*Method{
			CreateMethod(
				"isSuccess",
				[]string{"Boolean"},
				[]ast.Node{},
				func(this *Object, params []*Object, extra map[string]interface{}) interface{} {
					return this.Extra["isSuccess"]
				},
			),
		},
	)

	classMap := NewClassMap()
	saveResult := CreateClass(
		"SaveResult",
		[]*Method{},
		instanceMethods,
		NewMethodMap(),
	)
	classMap.Set("SaveResult", saveResult)
	nameSpaceStore.Set("Database", classMap)

	staticMethods := NewMethodMap()
	method := CreateMethod(
		"insert",
		[]string{"Database", "SaveResult"},
		[]ast.Node{objectTypeParameter}, // TODO: SObject or List<SObject>
		func(this *Object, params []*Object, extra map[string]interface{}) interface{} {
			sobj := params[0]
			record := &soapforce.SObject{}
			for k, v := range sobj.InstanceFields.All() {
				record.Fields[k] = v.Value()
			}
			client := NewSoapClient()
			rawSaveResults, err := client.Create([]*soapforce.SObject{record})
			if err != nil {
				panic(err)
			}
			retSaveResults := make([]*Object, len(rawSaveResults))
			for i, sr := range rawSaveResults {
				obj := CreateObject(saveResult)
				obj.Extra["isSuccess"] = NewBoolean(sr.Success)
				obj.Extra["id"] = NewString(sr.Id)
				obj.Extra["errors"] = sr.Errors
				retSaveResults[i] = obj
			}
			listObject := CreateObject(ListType)
			listObject.Extra["records"] = retSaveResults
			return listObject
		},
	)
	staticMethods.Set("insert", []*Method{method})
	method = CreateMethod(
		"setSavePoint",
		[]string{"Database", "SavePoint"},
		[]ast.Node{},
		func(this *Object, params []*Object, extra map[string]interface{}) interface{} {
			return Null
		},
	)
	staticMethods.Set("setSavePoint", []*Method{method})
	method = CreateMethod(
		"rollback",
		nil,
		[]ast.Node{objectTypeParameter}, // TODO: savepoint
		func(this *Object, params []*Object, extra map[string]interface{}) interface{} {
			return Null
		},
	)
	staticMethods.Set("rollback", []*Method{method})

	databaseClass := CreateClass(
		"Database",
		[]*Method{},
		nil,
		staticMethods,
	)
	primitiveClassMap.Set("Database", databaseClass)

}
