/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

type ICategory interface {
	GetDomainId() int

<<<<<<< HEAD
	GetValue() *ValueCategory
=======
	GetValue() ValueCategory
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

	SetValue(*ValueCategory) error

	Save() (int, error)

	// 获取子栏目的编号
	GetChildId() []int
}
