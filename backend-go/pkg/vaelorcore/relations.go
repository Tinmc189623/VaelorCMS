/*
 * VaelorCMS - 现代化内容管理系统
 *
 * Copyright © 2025-2026 Nexsteaduser. All rights reserved.
 *
 * 作者: Tinmc189623
 * 团队: Nexsteaduser
 *
 * 本程序是自由软件: 你可以重新分发和/或修改
 * 它在 GNU Affero 通用公共许可证的条款下发布,
 * 版本 3 或 (根据你的选择) 任何更高版本。
 *
 * 本程序是希望它有用,
 * 但没有任何保证; 甚至没有适销性或
 * 特定用途的默示保证。见
 * GNU Affero 通用公共许可证获取更多细节。
 *
 * 你应该收到 GNU Affero 通用公共许可证的副本
 * 与此程序一起。如果没有, 请见 <https://www.gnu.org/licenses/>.
 */

package vaelorcore

import (
	"reflect"
)

// RelationLoader 关系加载器
type RelationLoader struct {
	orm       *ORM
	relations map[string]Relation
}

// NewRelationLoader 创建新的关系加载器
// 参数: orm - ORM 实例
// 返回: 关系加载器实例
func NewRelationLoader(orm *ORM) *RelationLoader {
	return &RelationLoader{
		orm:       orm,
		relations: make(map[string]Relation),
	}
}

// RegisterRelation 注册关系
// 参数: name - 关系名称, relation - 关系定义
func (rl *RelationLoader) RegisterRelation(name string, relation Relation) {
	rl.relations[name] = relation
}

// LoadRelation 加载单个模型的关系
// 参数: model - 模型实例, relationName - 关系名称
// 返回: 可能的错误
func (rl *RelationLoader) LoadRelation(model Model, relationName string) error {
	relation, ok := rl.relations[relationName]
	if !ok {
		return ErrRelationNotFound
	}

	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}
	modelElem := modelValue.Elem()

	switch relation.Type {
	case HasOne, BelongsTo:
		return rl.loadHasOneOrBelongsTo(modelElem, relation)
	case HasMany:
		return rl.loadHasMany(modelElem, relation)
	case ManyToMany:
		return rl.loadManyToMany(model, modelElem, relation)
	default:
		return ErrInvalidRelation
	}
}

// LoadRelations 加载单个模型的多个关系
// 参数: model - 模型实例, relationNames - 关系名称列表
// 返回: 可能的错误
func (rl *RelationLoader) LoadRelations(model Model, relationNames ...string) error {
	for _, name := range relationNames {
		if err := rl.LoadRelation(model, name); err != nil {
			return err
		}
	}
	return nil
}

// LoadRelationForSlice 加载模型切片的关系
// 参数: models - 模型切片, relationName - 关系名称
// 返回: 可能的错误
func (rl *RelationLoader) LoadRelationForSlice(models interface{}, relationName string) error {
	relation, ok := rl.relations[relationName]
	if !ok {
		return ErrRelationNotFound
	}

	sliceValue := reflect.ValueOf(models)
	if sliceValue.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}
	sliceElem := sliceValue.Elem()
	if sliceElem.Kind() != reflect.Slice {
		return ErrMustBeSlice
	}

	if sliceElem.Len() == 0 {
		return nil
	}

	switch relation.Type {
	case HasOne, BelongsTo:
		return rl.loadHasOneOrBelongsToForSlice(sliceElem, relation)
	case HasMany:
		return rl.loadHasManyForSlice(sliceElem, relation)
	case ManyToMany:
		return rl.loadManyToManyForSlice(sliceElem, relation)
	default:
		return ErrInvalidRelation
	}
}

// LoadRelationsForSlice 加载模型切片的多个关系
// 参数: models - 模型切片, relationNames - 关系名称列表
// 返回: 可能的错误
func (rl *RelationLoader) LoadRelationsForSlice(models interface{}, relationNames ...string) error {
	for _, name := range relationNames {
		if err := rl.LoadRelationForSlice(models, name); err != nil {
			return err
		}
	}
	return nil
}

// loadHasOneOrBelongsTo 加载一对一或属于关系
func (rl *RelationLoader) loadHasOneOrBelongsTo(modelElem reflect.Value, relation Relation) error {
	localKey := relation.LocalKey
	if localKey == "" {
		localKey = "ID"
	}
	foreignKey := relation.ForeignKey
	if foreignKey == "" {
		if relation.Type == HasOne {
			foreignKey = toSnakeCase(modelElem.Type().Name()) + "_id"
		} else {
			foreignKey = toSnakeCase(reflect.TypeOf(relation.RelatedModel).Elem().Name()) + "_id"
		}
	}

	var keyValue interface{}
	if relation.Type == HasOne {
		keyValue = modelElem.FieldByName(localKey).Interface()
	} else {
		keyValue = modelElem.FieldByName(toCamelCase(foreignKey)).Interface()
	}

	relatedModel := reflect.New(reflect.TypeOf(relation.RelatedModel).Elem()).Interface().(Model)
	qb := NewQueryBuilder(relatedModel.TableName())

	if relation.Type == HasOne {
		qb.Where(foreignKey+" = ?", keyValue)
	} else {
		qb.Where("id = ?", keyValue)
	}

	found, err := rl.orm.FindFirst(relatedModel, qb)
	if err != nil {
		return err
	}

	field := modelElem.FieldByName(relation.Field)
	if !field.IsValid() || !field.CanSet() {
		return nil
	}

	if found {
		field.Set(reflect.ValueOf(relatedModel))
	}

	return nil
}

// loadHasOneOrBelongsToForSlice 加载模型切片的一对一或属于关系
func (rl *RelationLoader) loadHasOneOrBelongsToForSlice(sliceElem reflect.Value, relation Relation) error {
	localKey := relation.LocalKey
	if localKey == "" {
		localKey = "ID"
	}
	foreignKey := relation.ForeignKey
	if foreignKey == "" {
		elemType := sliceElem.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if relation.Type == HasOne {
			foreignKey = toSnakeCase(elemType.Name()) + "_id"
		} else {
			foreignKey = toSnakeCase(reflect.TypeOf(relation.RelatedModel).Elem().Name()) + "_id"
		}
	}

	keyMap := make(map[interface{}][]reflect.Value)
	var keyValues []interface{}

	for i := 0; i < sliceElem.Len(); i++ {
		elem := sliceElem.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		var keyValue interface{}
		if relation.Type == HasOne {
			keyValue = elem.FieldByName(localKey).Interface()
		} else {
			keyValue = elem.FieldByName(toCamelCase(foreignKey)).Interface()
		}

		keyMap[keyValue] = append(keyMap[keyValue], elem)
		keyValues = append(keyValues, keyValue)
	}

	if len(keyValues) == 0 {
		return nil
	}

	relatedType := reflect.TypeOf(relation.RelatedModel).Elem()
	relatedSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(relatedType))).Interface()

	qb := NewQueryBuilder(relation.RelatedModel.TableName())
	if relation.Type == HasOne {
		qb.WhereIn(foreignKey, keyValues)
	} else {
		qb.WhereIn("id", keyValues)
	}

	if err := rl.orm.FindAll(relatedSlice, qb); err != nil {
		return err
	}

	relatedSliceValue := reflect.ValueOf(relatedSlice).Elem()
	relatedMap := make(map[interface{}]reflect.Value)

	for i := 0; i < relatedSliceValue.Len(); i++ {
		relatedElem := relatedSliceValue.Index(i).Elem()
		var keyValue interface{}
		if relation.Type == HasOne {
			keyValue = relatedElem.FieldByName(toCamelCase(foreignKey)).Interface()
		} else {
			keyValue = relatedElem.FieldByName("ID").Interface()
		}
		relatedMap[keyValue] = relatedSliceValue.Index(i)
	}

	for keyValue, elems := range keyMap {
		related, ok := relatedMap[keyValue]
		for _, elem := range elems {
			field := elem.FieldByName(relation.Field)
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			if ok {
				field.Set(related)
			}
		}
	}

	return nil
}

// loadHasMany 加载一对多关系
func (rl *RelationLoader) loadHasMany(modelElem reflect.Value, relation Relation) error {
	localKey := relation.LocalKey
	if localKey == "" {
		localKey = "ID"
	}
	foreignKey := relation.ForeignKey
	if foreignKey == "" {
		foreignKey = toSnakeCase(modelElem.Type().Name()) + "_id"
	}

	keyValue := modelElem.FieldByName(localKey).Interface()

	relatedType := reflect.TypeOf(relation.RelatedModel).Elem()
	relatedSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(relatedType))).Interface()

	qb := NewQueryBuilder(relation.RelatedModel.TableName()).Where(foreignKey+" = ?", keyValue)
	if err := rl.orm.FindAll(relatedSlice, qb); err != nil {
		return err
	}

	field := modelElem.FieldByName(relation.Field)
	if !field.IsValid() || !field.CanSet() {
		return nil
	}

	field.Set(reflect.ValueOf(relatedSlice).Elem())

	return nil
}

// loadHasManyForSlice 加载模型切片的一对多关系
func (rl *RelationLoader) loadHasManyForSlice(sliceElem reflect.Value, relation Relation) error {
	localKey := relation.LocalKey
	if localKey == "" {
		localKey = "ID"
	}
	foreignKey := relation.ForeignKey
	if foreignKey == "" {
		elemType := sliceElem.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		foreignKey = toSnakeCase(elemType.Name()) + "_id"
	}

	keyMap := make(map[interface{}][]reflect.Value)
	var keyValues []interface{}

	for i := 0; i < sliceElem.Len(); i++ {
		elem := sliceElem.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		keyValue := elem.FieldByName(localKey).Interface()
		keyMap[keyValue] = append(keyMap[keyValue], elem)
		keyValues = append(keyValues, keyValue)
	}

	if len(keyValues) == 0 {
		return nil
	}

	relatedType := reflect.TypeOf(relation.RelatedModel).Elem()
	relatedSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(relatedType))).Interface()

	qb := NewQueryBuilder(relation.RelatedModel.TableName()).WhereIn(foreignKey, keyValues)
	if err := rl.orm.FindAll(relatedSlice, qb); err != nil {
		return err
	}

	relatedSliceValue := reflect.ValueOf(relatedSlice).Elem()
	relatedMap := make(map[interface{}][]reflect.Value)

	for i := 0; i < relatedSliceValue.Len(); i++ {
		relatedElem := relatedSliceValue.Index(i).Elem()
		keyValue := relatedElem.FieldByName(toCamelCase(foreignKey)).Interface()
		relatedMap[keyValue] = append(relatedMap[keyValue], relatedSliceValue.Index(i))
	}

	for keyValue, elems := range keyMap {
		relateds := relatedMap[keyValue]
		for _, elem := range elems {
			field := elem.FieldByName(relation.Field)
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			fieldSlice := reflect.MakeSlice(field.Type(), 0, len(relateds))
			fieldSlice = reflect.Append(fieldSlice, relateds...)
			field.Set(fieldSlice)
		}
	}

	return nil
}

// loadManyToMany 加载多对多关系
func (rl *RelationLoader) loadManyToMany(model Model, modelElem reflect.Value, relation Relation) error {
	pivotTable := relation.PivotTable
	if pivotTable == "" {
		return ErrInvalidRelation
	}
	pivotLocalKey := relation.PivotLocalKey
	if pivotLocalKey == "" {
		pivotLocalKey = toSnakeCase(modelElem.Type().Name()) + "_id"
	}
	pivotForeignKey := relation.PivotForeignKey
	if pivotForeignKey == "" {
		pivotForeignKey = toSnakeCase(reflect.TypeOf(relation.RelatedModel).Elem().Name()) + "_id"
	}

	var foreignIDs []int64
	qb := NewQueryBuilder(pivotTable).Select(pivotForeignKey).Where(pivotLocalKey+" = ?", model.GetID())
	if err := rl.orm.Pluck(model, pivotForeignKey, &foreignIDs, qb); err != nil {
		return err
	}

	if len(foreignIDs) == 0 {
		return nil
	}

	relatedType := reflect.TypeOf(relation.RelatedModel).Elem()
	relatedSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(relatedType))).Interface()

	qb = NewQueryBuilder(relation.RelatedModel.TableName()).WhereIn("id", interfaceSlice(foreignIDs))
	if err := rl.orm.FindAll(relatedSlice, qb); err != nil {
		return err
	}

	field := modelElem.FieldByName(relation.Field)
	if !field.IsValid() || !field.CanSet() {
		return nil
	}

	field.Set(reflect.ValueOf(relatedSlice).Elem())

	return nil
}

// loadManyToManyForSlice 加载模型切片的多对多关系
func (rl *RelationLoader) loadManyToManyForSlice(sliceElem reflect.Value, relation Relation) error {
	pivotTable := relation.PivotTable
	if pivotTable == "" {
		return ErrInvalidRelation
	}
	pivotLocalKey := relation.PivotLocalKey
	if pivotLocalKey == "" {
		elemType := sliceElem.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		pivotLocalKey = toSnakeCase(elemType.Name()) + "_id"
	}
	pivotForeignKey := relation.PivotForeignKey
	if pivotForeignKey == "" {
		pivotForeignKey = toSnakeCase(reflect.TypeOf(relation.RelatedModel).Elem().Name()) + "_id"
	}

	localIDMap := make(map[int64][]reflect.Value)
	var localIDs []int64

	for i := 0; i < sliceElem.Len(); i++ {
		elem := sliceElem.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		idValue := elem.FieldByName("ID").Interface().(int64)
		localIDMap[idValue] = append(localIDMap[idValue], elem)
		localIDs = append(localIDs, idValue)
	}

	if len(localIDs) == 0 {
		return nil
	}

	type PivotRow struct {
		LocalID   int64
		ForeignID int64
	}

	var pivotRows []*PivotRow
	qb := NewQueryBuilder(pivotTable).Select(pivotLocalKey, pivotForeignKey).WhereIn(pivotLocalKey, interfaceSlice(localIDs))

	query, args := qb.Build()
	rows, err := rl.orm.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var row PivotRow
		if err := rows.Scan(&row.LocalID, &row.ForeignID); err != nil {
			return err
		}
		pivotRows = append(pivotRows, &row)
	}

	if len(pivotRows) == 0 {
		return nil
	}

	var foreignIDs []int64
	foreignIDMap := make(map[int64][]int64)
	for _, row := range pivotRows {
		foreignIDs = append(foreignIDs, row.ForeignID)
		foreignIDMap[row.LocalID] = append(foreignIDMap[row.LocalID], row.ForeignID)
	}

	relatedType := reflect.TypeOf(relation.RelatedModel).Elem()
	relatedSlice := reflect.New(reflect.SliceOf(reflect.PtrTo(relatedType))).Interface()

	qb = NewQueryBuilder(relation.RelatedModel.TableName()).WhereIn("id", interfaceSlice(foreignIDs))
	if err := rl.orm.FindAll(relatedSlice, qb); err != nil {
		return err
	}

	relatedSliceValue := reflect.ValueOf(relatedSlice).Elem()
	relatedMap := make(map[int64]reflect.Value)
	for i := 0; i < relatedSliceValue.Len(); i++ {
		relatedElem := relatedSliceValue.Index(i).Elem()
		id := relatedElem.FieldByName("ID").Interface().(int64)
		relatedMap[id] = relatedSliceValue.Index(i)
	}

	for localID, elems := range localIDMap {
		foreignIDsForLocal := foreignIDMap[localID]
		var relateds []reflect.Value
		for _, fid := range foreignIDsForLocal {
			if related, ok := relatedMap[fid]; ok {
				relateds = append(relateds, related)
			}
		}

		for _, elem := range elems {
			field := elem.FieldByName(relation.Field)
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			fieldSlice := reflect.MakeSlice(field.Type(), 0, len(relateds))
			fieldSlice = reflect.Append(fieldSlice, relateds...)
			field.Set(fieldSlice)
		}
	}

	return nil
}

// interfaceSlice 将 int64 切片转换为 interface{} 切片
func interfaceSlice(slice []int64) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
