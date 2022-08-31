/**
    @auther: oreki
    @date: 2022/4/28
    @note: 图灵老祖保佑,永无BUG
**/

package controller

import (
	"context"
	"goods_srv/db"
	in "goods_srv/interface"
	"goods_srv/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌和轮播图
func (s *GoodsServer) BrandList(ctx context.Context, req *in.BrandFilterRequest) (*in.BrandListResponse, error) {
	var brandResponse []*in.BrandInfoResponse
	var brandListResponse in.BrandListResponse
	var brandList []*model.Brands
	var count int64
	db.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brandList)
	db.DB.Model(&model.Brands{}).Count(&count) // 获取总条数
	for _, brand := range brandList {
		brandResponse = append(brandResponse, &in.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Total = int32(count)
	brandListResponse.Data = brandResponse
	return &brandListResponse, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *in.BrandRequest) (*in.BrandInfoResponse, error) {
	//新建品牌
	if result := db.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	db.DB.Save(brand)

	return &in.BrandInfoResponse{Id: brand.ID}, nil
}
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *in.BrandRequest) (*emptypb.Empty, error) {
	if result := db.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *in.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := db.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	db.DB.Save(&brands)

	return &emptypb.Empty{}, nil
}
