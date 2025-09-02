package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"template-backend/internal/dto"
	"template-backend/internal/model"
	"template-backend/internal/repository"
)

type ResourceService interface {
	CreateResource(req *dto.CreateResourceRequest) (*dto.ResourceResponse, error)
	GetResourceByID(id int64) (*dto.ResourceResponse, error)
	UpdateResource(id int64, req *dto.UpdateResourceRequest) (*dto.ResourceResponse, error)
	DeleteResource(id int64) error
	ListResources(req *dto.ResourceQueryRequest) (*dto.PagedResponse, error)
	ResourcesTree(d *dto.ResourceQueryRequest) ([]dto.ResourceResponse, error)
}

type resourceService struct {
	resourceRepo repository.ResourceRepository
}

func NewResourceService(resourceRepo repository.ResourceRepository) ResourceService {
	return &resourceService{
		resourceRepo: resourceRepo,
	}
}

func (s *resourceService) CreateResource(req *dto.CreateResourceRequest) (*dto.ResourceResponse, error) {
	// 检查权限标识码是否已存在
	if s.resourceRepo.ExistsByPermissionCode(req.PermissionCode, 0) {
		return nil, errors.New("权限标识码已存在")
	}

	// 创建资源对象
	resource := &model.Resource{
		ResourceName:   req.ResourceName,
		PermissionCode: req.PermissionCode,
		Desc:           req.Desc,
		Type:           req.Type,
		ResourcePath:   req.ResourcePath,
		HTTPMethod:     req.HTTPMethod,
		ParentID:       req.ParentID,
		Sort:           req.Sort,
		RequiresAuth:   req.RequiresAuth,
		Remark:         req.Remark,
		CreatedBy:      req.CreatedBy,
		Status:         1, // 默认启用
	}

	if err := s.resourceRepo.Create(resource); err != nil {
		return nil, fmt.Errorf("创建资源失败: %w", err)
	}

	return s.modelToResponse(resource), nil
}

func (s *resourceService) GetResourceByID(id int64) (*dto.ResourceResponse, error) {
	resource, err := s.resourceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资源不存在")
		}
		return nil, fmt.Errorf("获取资源失败: %w", err)
	}

	return s.modelToResponse(resource), nil
}

func (s *resourceService) UpdateResource(id int64, req *dto.UpdateResourceRequest) (*dto.ResourceResponse, error) {
	// 检查资源是否存在
	existingResource, err := s.resourceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资源不存在")
		}
		return nil, fmt.Errorf("获取资源失败: %w", err)
	}

	// 如果要更新权限标识码，检查是否重复
	if req.PermissionCode != nil && *req.PermissionCode != existingResource.PermissionCode {
		if s.resourceRepo.ExistsByPermissionCode(*req.PermissionCode, id) {
			return nil, errors.New("权限标识码已存在")
		}
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.ResourceName != nil {
		updates["resource_name"] = *req.ResourceName
	}
	if req.PermissionCode != nil {
		updates["permission_code"] = *req.PermissionCode
	}
	if req.Desc != nil {
		updates["desc"] = *req.Desc
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.ResourcePath != nil {
		updates["resource_path"] = *req.ResourcePath
	}
	if req.HTTPMethod != nil {
		updates["http_method"] = *req.HTTPMethod
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.RequiresAuth != nil {
		updates["requires_auth"] = *req.RequiresAuth
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.UpdatedBy != nil {
		updates["updated_by"] = *req.UpdatedBy
	}

	// 执行更新
	if err := s.resourceRepo.Update(id, updates); err != nil {
		return nil, fmt.Errorf("更新资源失败: %w", err)
	}

	// 返回更新后的资源
	return s.GetResourceByID(id)
}

func (s *resourceService) DeleteResource(id int64) error {
	// 检查资源是否存在
	_, err := s.resourceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("资源不存在")
		}
		return fmt.Errorf("获取资源失败: %w", err)
	}

	// 执行删除
	if err := s.resourceRepo.Delete(id); err != nil {
		return fmt.Errorf("删除资源失败: %w", err)
	}

	return nil
}

func (s *resourceService) ListResources(req *dto.ResourceQueryRequest) (*dto.PagedResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	query := &repository.ResourceQuery{
		ID:             req.ID,
		ResourceName:   req.ResourceName,
		PermissionCode: req.PermissionCode,
		Type:           req.Type,
		ResourcePath:   req.ResourcePath,
		HTTPMethod:     req.HTTPMethod,
		ParentID:       req.ParentID,
		Status:         req.Status,
		RequiresAuth:   req.RequiresAuth,
		Page:           req.Page,
		PageSize:       req.PageSize,
	}

	resources, total, err := s.resourceRepo.List(query)
	if err != nil {
		return nil, fmt.Errorf("查询资源失败: %w", err)
	}

	// 转换为响应格式
	data := make([]dto.ResourceResponse, len(resources))
	for i, resource := range resources {
		data[i] = *s.modelToResponse(&resource)
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PagedResponse{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *resourceService) ResourcesTree(req *dto.ResourceQueryRequest) ([]dto.ResourceResponse, error) {
	query := &repository.ResourceQuery{
		ID:             req.ID,
		ResourceName:   req.ResourceName,
		PermissionCode: req.PermissionCode,
		Type:           req.Type,
		ResourcePath:   req.ResourcePath,
		HTTPMethod:     req.HTTPMethod,
		ParentID:       req.ParentID,
		Status:         req.Status,
		RequiresAuth:   req.RequiresAuth,
		Page:           req.Page,
		PageSize:       req.PageSize,
	}

	resources, err := s.resourceRepo.QueryAll(query)
	if err != nil {
		return nil, fmt.Errorf("查询资源失败: %w", err)
	}

	// 转换为响应格式
	data := make([]dto.ResourceResponse, len(resources))
	for i, resource := range resources {
		data[i] = *s.modelToResponse(&resource)
	}

	if req.ParentID != nil {
		return data, nil
	}

	var roots []dto.ResourceResponse
	menuMap := make(map[int64]*dto.ResourceResponse)
	for _, m := range roots {
		menuMap[m.ID] = &m
	}
	for _, m := range data {
		if m.ParentID != nil {
			if parent, ok := menuMap[*m.ParentID]; ok {
				parent.Children = append(parent.Children, m)
			}
		} else {
			roots = append(roots, m)
		}
	}

	return roots, nil
}

func (s *resourceService) modelToResponse(resource *model.Resource) *dto.ResourceResponse {
	return &dto.ResourceResponse{
		ID:             resource.ID,
		ResourceName:   resource.ResourceName,
		PermissionCode: resource.PermissionCode,
		Desc:           resource.Desc,
		Type:           resource.Type,
		ResourcePath:   resource.ResourcePath,
		HTTPMethod:     resource.HTTPMethod,
		ParentID:       resource.ParentID,
		Sort:           resource.Sort,
		Status:         resource.Status,
		RequiresAuth:   resource.RequiresAuth,
		Remark:         resource.Remark,
		CreatedBy:      resource.CreatedBy,
		CreatedAt:      resource.CreatedAt,
		UpdatedBy:      resource.UpdatedBy,
		UpdatedAt:      resource.UpdatedAt,
		HasChildren:    true,
	}
}
