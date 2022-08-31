/**
    @auther: oreki
    @date: 2022/5/12
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

// 轮播图
func (s *GoodsServer) BannerList(ctx context.Context, req *emptypb.Empty) (*in.BannerListResponse, error) {
	bannerListResponse := in.BannerListResponse{}

	var banners []model.Banner
	result := db.DB.Find(&banners)
	bannerListResponse.Total = int32(result.RowsAffected)

	var bannerReponses []*in.BannerResponse
	for _, banner := range banners {
		bannerReponses = append(bannerReponses, &in.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}

	bannerListResponse.Data = bannerReponses

	return &bannerListResponse, nil
}

func (s *GoodsServer) CreateBanner(ctx context.Context, req *in.BannerRequest) (*in.BannerResponse, error) {
	banner := model.Banner{}

	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	db.DB.Save(&banner)

	return &in.BannerResponse{Id: banner.ID}, nil
}

func (s *GoodsServer) DeleteBanner(ctx context.Context, req *in.BannerRequest) (*emptypb.Empty, error) {
	if result := db.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBanner(ctx context.Context, req *in.BannerRequest) (*emptypb.Empty, error) {
	var banner model.Banner

	if result := db.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	db.DB.Save(&banner)

	return &emptypb.Empty{}, nil
}
