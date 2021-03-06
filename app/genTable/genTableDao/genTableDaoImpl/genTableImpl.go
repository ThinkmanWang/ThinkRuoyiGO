package genTableDaoImpl

import (
	"baize/app/common/mysql"
	"baize/app/constant/constants"
	"baize/app/genTable/genTableModels"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var genTableDaoImpl = &genTableDao{db: mysql.GetMysqlDb()}

type genTableDao struct {
	db **sqlx.DB
}

func GetGenTableDao() *genTableDao {
	return genTableDaoImpl
}
func (genTableDao *genTableDao) getDb() *sqlx.DB {
	return *genTableDao.db
}

func (genTableDao *genTableDao) SelectGenTableList(GenTable *genTableModels.GenTableDQL) (list []*genTableModels.GenTableVo, total *int64) {
	var selectSql = `select table_id, table_name, table_comment, sub_table_name, sub_table_fk_name, class_name, tpl_category, package_name, module_name, business_name, function_name, function_author, gen_type, gen_path, options, create_by, create_time, update_by, update_time, remark `
	var fromSql = ` from gen_table`
	whereSql := ``
	if GenTable.TableName != "" {
		whereSql += " AND lower(table_name) like lower(concat('%', :table_name, '%'))"
	}
	if GenTable.TableComment != "" {
		whereSql += " AND lower(table_comment) like lower(concat('%', :table_comment, '%'))"
	}
	if GenTable.BeginTime != "" {
		whereSql += " AND date_format(create_time,'%y%m%d') &gt;= date_format(:begin_time,'%y%m%d')"
	}
	if GenTable.EndTime != "" {
		whereSql += " date_format(create_time,'%y%m%d') &lt;= date_format(:end_time,'%y%m%d')"
	}

	if whereSql != "" {
		whereSql = " where " + whereSql[4:]
	}
	countSql := constants.MysqlCount + fromSql + whereSql

	countRow, err := genTableDao.getDb().NamedQuery(countSql, GenTable)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()
	list = make([]*genTableModels.GenTableVo, 0, GenTable.Size)
	if *total > GenTable.Offset {
		if GenTable.Limit != "" {
			whereSql += GenTable.Limit
		}
		listRows, err := genTableDao.getDb().NamedQuery(selectSql+fromSql+whereSql, GenTable)
		if err != nil {
			panic(err)
		}
		for listRows.Next() {
			postVo := new(genTableModels.GenTableVo)
			err := listRows.StructScan(postVo)
			if err != nil {
				panic(err)
			}
			list = append(list, postVo)
		}
		defer listRows.Close()
	}
	return
}
func (genTableDao *genTableDao) SelectDbTableList(GenTable *genTableModels.GenTableDQL) (list []*genTableModels.DBTableVo, total *int64) {
	var selectSql = `select table_name , table_comment, create_time, update_time `
	var fromSql = ` from information_schema.tables`
	whereSql := ` where table_schema = (select database())
		AND table_name NOT LIKE 'qrtz_%' AND table_name NOT LIKE 'gen_%'
		AND table_name NOT IN (select table_name from gen_table)`
	if GenTable.TableName != "" {
		whereSql += " AND lower(table_name) like lower(concat('%', :table_name, '%'))"
	}
	if GenTable.TableComment != "" {
		whereSql += " AND lower(table_comment) like lower(concat('%', :table_comment, '%'))"
	}
	if GenTable.BeginTime != "" {
		whereSql += " AND date_format(create_time,'%y%m%d') &gt;= date_format(:begin_time,'%y%m%d')"
	}
	if GenTable.EndTime != "" {
		whereSql += " date_format(create_time,'%y%m%d') &lt;= date_format(:end_time,'%y%m%d')"
	}

	countSql := constants.MysqlCount + fromSql + whereSql

	countRow, err := genTableDao.getDb().NamedQuery(countSql, GenTable)
	if err != nil {
		panic(err)
	}
	total = new(int64)
	if countRow.Next() {
		countRow.Scan(total)
	}
	defer countRow.Close()
	list = make([]*genTableModels.DBTableVo, 0, GenTable.Size)
	if *total > GenTable.Offset {
		if GenTable.Limit != "" {
			whereSql += GenTable.Limit
		}
		listRows, err := genTableDao.getDb().NamedQuery(selectSql+fromSql+whereSql, GenTable)
		if err != nil {
			panic(err)
		}
		for listRows.Next() {
			dbTable := new(genTableModels.DBTableVo)
			err := listRows.StructScan(dbTable)
			if err != nil {
				panic(err)
			}
			list = append(list, dbTable)
		}
		defer listRows.Close()
	}
	return
}

func (genTableDao *genTableDao) SelectDbTableListByNames(tableNames []string) (list []*genTableModels.DBTableVo) {
	query, i, err := sqlx.In("select table_name, table_comment, create_time, update_time from information_schema.tables where table_name NOT LIKE 'gen_%' and table_schema = (select database()) and table_name in  (?)", tableNames)
	if err != nil {
		panic(err)
	}
	list = make([]*genTableModels.DBTableVo, 0, 0)
	err = genTableDao.getDb().Select(&list, query, i...)
	if err != nil {
		panic(err)
	}
	return
}

func (genTableDao *genTableDao) SelectGenTableById(id int64) (genTable *genTableModels.GenTableVo) {
	genTable = new(genTableModels.GenTableVo)
	err := genTableDao.getDb().Get(genTable, `SELECT
       table_id, table_name, table_comment, sub_table_name,sub_table_fk_name, class_name, 
      tpl_category, package_name,module_name, business_name,function_name, function_author,gen_type,gen_path, options, remark
		FROM gen_table 
		where table_id = ?`, id)
	if err != nil {
		panic(err)
	}
	return
}
func (genTableDao *genTableDao) SelectGenTableByName(name string) (genTable *genTableModels.GenTableVo) {
	genTable = new(genTableModels.GenTableVo)
	err := genTableDao.getDb().Get(genTable, `SELECT t.table_id, t.table_name, t.table_comment, t.sub_table_name, t.sub_table_fk_name, t.class_name, t.tpl_category, t.package_name, t.module_name, t.business_name, t.function_name, t.function_author, t.gen_type, t.gen_path, t.options, t.remark
		FROM gen_table t
		where t.table_name = ? `, name)
	if err != nil {
		panic(err)
	}
	return
}
func (genTableDao *genTableDao) SelectGenTableAll() (list []*genTableModels.GenTableVo) {
	list = make([]*genTableModels.GenTableVo, 0, 0)
	err := genTableDao.getDb().Select(&list, `SELECT t.table_id, t.table_name, t.table_comment, t.sub_table_name, t.sub_table_fk_name, t.class_name, t.tpl_category, t.package_name, t.module_name, t.business_name, t.function_name, t.function_author, t.gen_type, t.gen_path, t.options, t.remark
		FROM gen_table t`)
	if err != nil {
		panic(err)
	}
	return
}

func (genTableDao *genTableDao) BatchInsertGenTable(genTables []*genTableModels.GenTableDML) {

	_, err := genTableDao.getDb().NamedExec(`insert into gen_table(table_id,table_name,table_comment,class_name,tpl_category,package_name,module_name,business_name,function_name,function_author,gen_type,gen_path,create_by,create_time,update_by,update_time,remark)
							values(:table_id,:table_name,:table_comment,:class_name,:tpl_category,:package_name,:module_name,:business_name,:function_name,:function_author,:gen_type,:gen_path,:create_by,now(),:update_by,now(),:remark)`,
		genTables)
	if err != nil {
		panic(err)
	}

}

func (genTableDao *genTableDao) InsertGenTable(genTable *genTableModels.GenTableDML) {
	insertSQL := `insert into gen_table(table_id,table_name,table_comment,class_name,tpl_category,package_name,module_name,business_name,function_name,function_author,gen_type,gen_path,create_by,create_time,update_by,update_time %s)
					values(:table_id,:table_name,:table_comment,:class_name,:tpl_category,:package_name,:module_name,:business_name,:function_name,:function_author,:gen_type,:gen_path,:create_by,now(),:update_by,now() %s)`
	key := ""
	value := ""

	if genTable.Remark != "" {
		key += ",remark"
		value += ",:remark"
	}

	insertStr := fmt.Sprintf(insertSQL, key, value)
	_, err := genTableDao.getDb().NamedExec(insertStr, genTable)
	if err != nil {
		panic(err)
	}
	return
}

func (genTableDao *genTableDao) UpdateGenTable(genTable *genTableModels.GenTableDML) {
	updateSQL := `update gen_table set update_time = now() , update_by = :update_by`
	if genTable.TableName != "" {
		updateSQL += ",table_name = :table_name"
	}
	if genTable.TableComment != "" {
		updateSQL += ",table_comment = :table_comment"
	}
	if genTable.SubTableName != "" {
		updateSQL += ",sub_table_name = :sub_table_name"
	}
	if genTable.SubTableFkName != "" {
		updateSQL += ",sub_table_fk_name = :sub_table_fk_name"
	}
	if genTable.ClassName != "" {
		updateSQL += ",class_name = :class_name"
	}
	if genTable.FunctionAuthor != "" {
		updateSQL += ",function_author = :function_author"
	}
	if genTable.GenType != "" {
		updateSQL += ",gen_type = :gen_type"
	}
	if genTable.GenPath != "" {
		updateSQL += ",gen_path = :gen_path"
	}
	if genTable.TplCategory != "" {
		updateSQL += ",tpl_category = :tpl_category"
	}
	if genTable.PackageName != "" {
		updateSQL += ",package_name = :package_name"
	}
	if genTable.ModuleName != "" {
		updateSQL += ",module_name = :module_name"
	}
	if genTable.BusinessName != "" {
		updateSQL += ",business_name = :business_name"
	}
	if genTable.FunctionName != "" {
		updateSQL += ",function_name = :function_name"
	}
	if genTable.Options != "" {
		updateSQL += ",options = :options"
	}
	if genTable.Remark != "" {
		updateSQL += ",remark = :remark"
	}

	updateSQL += " where table_id = :table_id"

	_, err := genTableDao.getDb().NamedExec(updateSQL, genTable)
	if err != nil {
		panic(err)
	}
	return
}

func (genTableDao *genTableDao) DeleteGenTableByIds(ids []int64) {
	query, i, err := sqlx.In(" delete from gen_table where table_id in(?)", ids)
	if err != nil {
		panic(err)
	}
	_, err = genTableDao.getDb().Exec(query, i...)
	if err != nil {
		panic(err)
	}

}
