package genTableModels

import (
	"baize/app/common/commonModels"
	genUtils "baize/app/genTable/utils"
	"baize/app/utils/snowflake"
	"baize/app/utils/stringUtils"
	"strconv"
	"strings"
)

type GenTableColumnDML struct {
	ColumnId      int64  `json:"columnId,string" db:"column_id"`
	TableId       int64  `json:"tableId,string" db:"table_id"`
	ColumnName    string `json:"columnName" db:"column_name"`
	ColumnComment string `json:"columnComment" db:"column_comment"`
	ColumnType    string `json:"columnType" db:"column_type"`
	GoType        string `json:"goType" db:"go_type"`
	GoField       string `json:"goField" db:"go_field"`
	HtmlField     string `json:"htmlField" db:"html_field"`
	IsPk          string `json:"isPk" db:"is_pk"`
	IsRequired    string `json:"isRequired" db:"is_required"`
	IsInsert      string `json:"isInsert" db:"is_insert"`
	IsEdit        string `json:"isEdit" db:"is_edit"`
	IsList        string `json:"isList" db:"is_list"`
	IsQuery       string `json:"isQuery" db:"is_query"`
	QueryType     string `json:"queryType" db:"query_type"`
	HtmlType      string `json:"htmlType" db:"html_type"`
	DictType      string `json:"dictType" db:"dict_type"`
	Sort          int32  `json:"sort" db:"sort"`
	commonModels.BaseEntityDML
}

type GenTableColumnVo struct {
	ColumnId      int64   `json:"columnId,string" db:"column_id"`
	TableId       int64   `json:"tableId,string" db:"table_id"`
	ColumnName    string  `json:"columnName" db:"column_name"`
	ColumnComment string  `json:"columnComment" db:"column_comment"`
	ColumnType    string  `json:"columnType" db:"column_type"`
	GoType        string  `json:"goType" db:"go_type"`
	GoField       string  `json:"goField" db:"go_field"`
	HtmlField     string  `json:"htmlField" db:"html_field"`
	IsPk          string  `json:"isPk" db:"is_pk"`
	IsRequired    string  `json:"isRequired" db:"is_required"`
	IsInsert      string  `json:"isInsert" db:"is_insert"`
	IsEdit        string  `json:"isEdit" db:"is_edit"`
	IsList        string  `json:"isList" db:"is_list"`
	IsQuery       string  `json:"isQuery" db:"is_query"`
	QueryType     string  `json:"queryType" db:"query_type"`
	HtmlType      string  `json:"htmlType" db:"html_type"`
	DictType      *string `json:"dictType" db:"dict_type"`
	Sort          int32   `json:"remark" db:"sort"`

	commonModels.BaseEntity
}

type InformationSchemaColumn struct {
	ColumnName    string `db:"COLUMN_NAME"`
	ColumnComment string `db:"COLUMN_COMMENT"`
	ColumnType    string `db:"COLUMN_TYPE"`
	IsPk          string `db:"is_pk"`
	IsRequired    string `db:"is_required"`
	Sort          int32  `db:"sort"`
}

func GetGenTableColumnDML(column *InformationSchemaColumn, tableId int64, userName string) *GenTableColumnDML {
	genTableColumn := new(GenTableColumnDML)
	dataType := genUtils.GetDbType(column.ColumnType)
	columnName := column.ColumnName
	genTableColumn.ColumnId = snowflake.GenID()
	genTableColumn.ColumnName = column.ColumnName
	genTableColumn.IsPk = column.IsPk
	genTableColumn.Sort = column.Sort
	genTableColumn.ColumnComment = column.ColumnComment
	genTableColumn.ColumnType = column.ColumnType
	genTableColumn.TableId = tableId
	genTableColumn.CreateBy = userName
	genTableColumn.UpdateBy = userName
	//???????????????
	genTableColumn.GoField = stringUtils.ConvertToBigCamelCase(columnName)
	genTableColumn.HtmlField = stringUtils.ConvertToLittleCamelCase(columnName)

	if genUtils.IsStringObject(dataType) {
		//????????????????????????
		genTableColumn.GoType = "string"
		if strings.EqualFold(dataType, "text") || strings.EqualFold(dataType, "tinytext") || strings.EqualFold(dataType, "mediumtext") || strings.EqualFold(dataType, "longtext") {
			genTableColumn.HtmlType = "textarea"
		} else {
			columnLength := genUtils.GetColumnLength(column.ColumnType)
			if columnLength >= 500 {
				genTableColumn.HtmlType = "textarea"
			} else {
				genTableColumn.HtmlType = "input"
			}
		}
	} else if genUtils.IsTimeObject(dataType) {
		//?????????????????????
		genTableColumn.GoType = "Time"
		genTableColumn.HtmlType = "datetime"
	} else if genUtils.IsNumberObject(dataType) {
		//?????????????????????
		genTableColumn.HtmlType = "input"
		// ??????????????????
		tmp := genTableColumn.ColumnType
		if tmp == "float" || tmp == "double" {
			genTableColumn.GoType = "float64"
		} else {
			start := strings.Index(tmp, "(")
			end := strings.Index(tmp, ")")
			if end < 0 {
				genTableColumn.GoType = "int64"
			} else {
				result := tmp[start+1 : end]
				arr := strings.Split(result, ",")
				i0, _ := strconv.Atoi(arr[0])
				i1, _ := strconv.Atoi(arr[1])
				if len(arr) == 2 && i0 > 0 {
					genTableColumn.GoType = "float64"
				} else if len(arr) == 1 && i1 <= 10 {
					genTableColumn.GoType = "int"
				} else {
					genTableColumn.GoType = "int64"
				}
			}
		}

	}
	//????????????
	if columnName == "create_by" || columnName == "create_time" || columnName == "update_by" || columnName == "update_time" {
		genTableColumn.IsRequired = "0"
		genTableColumn.IsInsert = "0"
	} else {
		genTableColumn.IsRequired = "0"
		genTableColumn.IsInsert = "1"
		if strings.Index(columnName, "name") >= 0 || strings.Index(columnName, "status") >= 0 {
			genTableColumn.IsRequired = "1"
		}
	}

	// ????????????
	if genUtils.IsNotEdit(columnName) {
		if column.IsPk == "1" {
			genTableColumn.IsEdit = "0"
		} else {
			genTableColumn.IsEdit = "1"
		}
	} else {
		genTableColumn.IsEdit = "0"
	}
	// ????????????
	if genUtils.IsNotList(columnName) {
		genTableColumn.IsList = "1"
	} else {
		genTableColumn.IsList = "0"
	}
	// ????????????
	if genUtils.IsNotQuery(columnName) {
		genTableColumn.IsQuery = "1"
	} else {
		genTableColumn.IsQuery = "0"
	}

	// ??????????????????
	if genUtils.CheckNameColumn(columnName) {
		genTableColumn.QueryType = "LIKE"
	} else {
		genTableColumn.QueryType = "EQ"
	}

	// ???????????????????????????
	if genUtils.CheckStatusColumn(columnName) {
		genTableColumn.HtmlType = "radio"
	} else if genUtils.CheckTypeColumn(columnName) || genUtils.CheckSexColumn(columnName) {
		// ??????&???????????????????????????
		genTableColumn.HtmlType = "select"
	}
	return genTableColumn
}
