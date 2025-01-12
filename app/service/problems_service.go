package service

import (
	"fmt"
	"oj-back/app/db"
	"oj-back/app/model"
)

type ProblemService struct{}

var ProblemServiceApp = new(ProblemService)

// GetAllProblems 查询所有题目并分页
func (ps *ProblemService) GetAllProblems(page int, size int, keyword string) (*model.Page[model.Problem], error) {
	var problems []model.Problem
	var total int64
	query := db.DB.Model(&model.Problem{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取题目总数失败: %w", err)
	}

	// 获取题目数据
	err := query.Offset((page - 1) * size).Limit(size).Find(&problems).Error
	if err != nil {
		return nil, fmt.Errorf("获取题目列表失败: %w", err)
	}

	// 返回分页结果
	return &model.Page[model.Problem]{
		Total: total,
		Data:  problems,
	}, nil
}

// GetProblemByID 查询指定 ID 的题目详情
func (ps *ProblemService) GetProblemByID(id int) (*model.Problem, error) {
	var problem model.Problem
	err := db.DB.Model(&model.Problem{}).Where("id = ?", id).First(&problem).Error
	if err != nil {
		return nil, fmt.Errorf("获取题目详情失败: %w", err)
	}

	// 返回题目详情
	return &problem, nil
}

func (ps *ProblemService) GetProblemTestCase(problemID uint64) ([]model.Case, error) {
	var record model.TestCase

	err := db.DB.Where("problem_id = ?", problemID).First(&record).Error
	if err != nil {
		return nil, fmt.Errorf("查询测试用例失败: %v", err)
	}
	cases := record.Cases

	return cases, nil
}
