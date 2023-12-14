package services

import (
	"context"
	"go-grpc/cmd/helpers"
	"go-grpc/pb"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	*gorm.DB
}

func (p *ProductService) GetProducts(ctx context.Context, param *pb.Page) (*pb.Products, error) {
	var page int64 = 1

	if param.GetPage() != 0 {
		page = param.GetPage()
	}

	var pagination pb.Pagination
	var products []*pb.Product

	sql := p.DB.Table("products as p").Joins("LEFT JOIN categories as c ON c.id = p.category_id").Select("p.id", "p.name", "p.price", "p.stock", "c.id as category_id", "c.name as category_name")

	offest, limit := helpers.Pagination(sql, page, &pagination)
	rows, err := sql.Offset(int(offest)).Limit(int(limit)).Rows()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var product pb.Product
		var category pb.Category

		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
			log.Fatal("Gagal mengambil row data %v", err.Error())
		}

		product.Category = &category
		products = append(products, &product)
	}

	response := &pb.Products{
		Pagination: &pagination,
		Data:       products,
	}

	return response, nil
}

func (p *ProductService) GetProduct(ctx context.Context, id *pb.Id) (*pb.Product, error) {
	row := p.DB.Table("products as p").
		Joins("LEFT JOIN categories as c ON c.id = p.category_id").
		Select("p.id", "p.name", "p.price", "p.stock", "c.id as category_id", "c.name as category_name").
		Where("p.id = ?", id.GetId()).
		Row()

	var product pb.Product
	var category pb.Category

	if err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	product.Category = &category

	return &product, nil
}

func (p *ProductService) CreateProduct(context context.Context, productData *pb.Product) (*pb.Id, error) {
	var Response pb.Id

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := pb.Category{
			Id:   0,
			Name: productData.GetCategory().GetName(),
		}

		if err := tx.Table("categories").
			Where("name = ?", category.GetName()).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			Id         uint64
			Name       string
			Price      uint64
			Stock      uint64
			CategoryId uint32
		}{
			Id:         productData.GetId(),
			Name:       productData.GetName(),
			Price:      uint64(productData.GetPrice()),
			Stock:      uint64(productData.GetStock()),
			CategoryId: category.GetId(),
		}

		if err := tx.Table("products").Create(&product).Error; err != nil {
			return err
		}

		Response.Id = product.Id

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

func (p *ProductService) UpdateProduct(context context.Context, productData *pb.Product) (*pb.Status, error) {
	var Response pb.Status

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := pb.Category{
			Id:   0,
			Name: productData.GetCategory().GetName(),
		}

		if err := tx.Table("categories").
			Where("name = ?", category.GetName()).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			Id         uint64
			Name       string
			Price      uint64
			Stock      uint64
			CategoryId uint32
		}{
			Id:         productData.GetId(),
			Name:       productData.GetName(),
			Price:      uint64(productData.GetPrice()),
			Stock:      uint64(productData.GetStock()),
			CategoryId: category.GetId(),
		}

		if err := tx.Table("products").Where("id = ?", product.Id).Updates(&product).Error; err != nil {
			return err
		}

		Response.Status = 1

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

func (p *ProductService) DeleteProduct(context context.Context, id *pb.Id) (*pb.Status, error) {
	var response pb.Status

	if err := p.DB.Table("products").Where("id = ?", id.GetId()).Delete(nil).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Status = 1

	return &response, nil
}
