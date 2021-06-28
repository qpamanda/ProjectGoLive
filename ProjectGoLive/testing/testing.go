package testing

// Author: Ahmad Bahrudin
import (
	"ProjectGoLive/database"
	"fmt"
)

func TestCatInsert() {
	database.Cat.Insert(4, "Category TESTING", "admin")
	database.Cat.Insert(5, "Category TESTING", "admin")
	database.Cat.Insert(6, "Category TESTING", "admin")
}

func TestCatUpdate() {
	database.Cat.Update(4, "Category TESTING 1", "admin")
	database.Cat.Update(5, "Category TESTING 1", "admin")
	database.Cat.Update(6, "Category TESTING 1", "admin")
}

func TestCatGet() {
	cat1, _ := database.Cat.Get(1)
	cat2, _ := database.Cat.Get(2)
	cat3, _ := database.Cat.Get(3)

	fmt.Println(cat1.CategoryID, cat1.Category, cat1.CreatedBy, cat1.Created_dt, cat1.LastModifiedBy, cat1.LastModified_dt)
	fmt.Println(cat2.CategoryID, cat2.Category, cat2.CreatedBy, cat2.Created_dt, cat2.LastModifiedBy, cat2.LastModified_dt)
	fmt.Println(cat3.CategoryID, cat3.Category, cat3.CreatedBy, cat3.Created_dt, cat3.LastModifiedBy, cat3.LastModified_dt)
	fmt.Println()
}

func TestCatGetAll() {
	cat, _ := database.Cat.GetAll()

	for i := 1; i <= len(cat); i++ {
		fmt.Println(cat[i].CategoryID, cat[i].Category, cat[i].CreatedBy, cat[i].Created_dt, cat[i].LastModifiedBy, cat[i].LastModified_dt)
	}
	fmt.Println()
}

func TestCatDelete() {
	database.Cat.Delete(4)
	database.Cat.Delete(5)
	database.Cat.Delete(6)
}

func TestMemTInsert() {
	database.MemT.Insert(4, "MemberType TESTING", "admin")
	database.MemT.Insert(5, "MemberType TESTING", "admin")
	database.MemT.Insert(6, "MemberType TESTING", "admin")
}

func TestMemTUpdate() {
	database.MemT.Update(4, "MemberType TESTING 1", "admin")
	database.MemT.Update(5, "MemberType TESTING 1", "admin")
	database.MemT.Update(6, "MemberType TESTING 1", "admin")
}

func TestMemTGet() {
	memT1, _ := database.MemT.Get(1)
	memT2, _ := database.MemT.Get(2)
	memT3, _ := database.MemT.Get(3)

	fmt.Println(memT1.MemberTypeID, memT1.MemberType, memT1.CreatedBy, memT1.Created_dt, memT1.LastModifiedBy, memT1.LastModified_dt)
	fmt.Println(memT2.MemberTypeID, memT2.MemberType, memT2.CreatedBy, memT2.Created_dt, memT2.LastModifiedBy, memT2.LastModified_dt)
	fmt.Println(memT3.MemberTypeID, memT3.MemberType, memT3.CreatedBy, memT3.Created_dt, memT3.LastModifiedBy, memT3.LastModified_dt)
	fmt.Println()
}

func TestMemTGetAll() {
	memT, _ := database.MemT.GetAll()

	for i := 1; i <= len(memT); i++ {
		fmt.Println(memT[i].MemberTypeID, memT[i].MemberType, memT[i].CreatedBy, memT[i].Created_dt, memT[i].LastModifiedBy, memT[i].LastModified_dt)
	}
	fmt.Println()
}

func TestMemTDelete() {
	database.MemT.Delete(4)
	database.MemT.Delete(5)
	database.MemT.Delete(6)
}

func TestReqSInsert() {
	database.ReqS.Insert(3, "RequestStatus TESTING", "admin")
	database.ReqS.Insert(4, "RequestStatus TESTING", "admin")
	database.ReqS.Insert(5, "RequestStatus TESTING", "admin")
}

func TestReqSUpdate() {
	database.ReqS.Update(3, "RequestStatus TESTING 1", "admin")
	database.ReqS.Update(4, "RequestStatus TESTING 1", "admin")
	database.ReqS.Update(5, "RequestStatus TESTING 1", "admin")
}

func TestReqSGet() {
	reqS1, _ := database.ReqS.Get(0)
	reqS2, _ := database.ReqS.Get(1)
	reqS3, _ := database.ReqS.Get(2)

	fmt.Println(reqS1.StatusCode, reqS1.Status, reqS1.CreatedBy, reqS1.Created_dt, reqS1.LastModifiedBy, reqS1.LastModified_dt)
	fmt.Println(reqS2.StatusCode, reqS2.Status, reqS2.CreatedBy, reqS2.Created_dt, reqS2.LastModifiedBy, reqS2.LastModified_dt)
	fmt.Println(reqS3.StatusCode, reqS3.Status, reqS3.CreatedBy, reqS3.Created_dt, reqS3.LastModifiedBy, reqS3.LastModified_dt)
	fmt.Println()
}

func TestReqSGetAll() {
	reqS, _ := database.ReqS.GetAll()

	for i := 0; i < len(reqS); i++ {
		fmt.Println(reqS[i].StatusCode, reqS[i].Status, reqS[i].CreatedBy, reqS[i].Created_dt, reqS[i].LastModifiedBy, reqS[i].LastModified_dt)
	}
	fmt.Println()
}

func TestReqSDelete() {
	database.ReqS.Delete(3)
	database.ReqS.Delete(4)
	database.ReqS.Delete(5)
}
