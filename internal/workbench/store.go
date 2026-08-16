package workbench

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("article not found")

type Store struct {
	mu       sync.RWMutex
	articles map[string]Article
	order    []string
}

func NewFixtureStore() *Store {
	fixtures := []Article{
		{
			ID:           "news-library",
			Section:      SectionNews,
			Title:        "新图书馆试运行首周开放",
			Summary:      "新馆阅览区、研讨室与自助借还系统面向师生开放试运行。",
			Body:         "本周一，学校新图书馆进入试运行阶段。首批开放区域包括一至三层阅览区、十二间小组研讨室和自助借还区。\n\n图书馆老师介绍，试运行期间将根据师生反馈调整座位预约规则，并逐步延长晚间开放时间。学生可通过校园卡进入新馆，馆藏迁移工作仍在有序进行。",
			Author:       "林小满",
			Status:       StatusPendingReview,
			Edition:      "第 128 期",
			UpdatedLabel: "今天 10:24",
		},
		{
			ID:           "news-lab",
			Section:      SectionNews,
			Title:        "创客实验室新增开放时段",
			Summary:      "创客实验室从本月起增加周末预约时段。",
			Body:         "为满足学生社团和项目小组的制作需求，创客实验室从本月起增加周六上午和周日下午两个开放时段。使用设备前仍需完成安全培训。",
			Author:       "周嘉树",
			Status:       StatusDraft,
			Edition:      "第 128 期",
			UpdatedLabel: "昨天 16:40",
		},
		{
			ID:           "interview-alumni",
			Section:      SectionInterview,
			Title:        "校友访谈：在田野里寻找答案",
			Summary:      "青年地理学者陈乔分享十年野外调查经历。",
			Body:         "从校园地理社第一次野外考察，到独立主持高原生态调查，陈乔始终把现场观察看作研究的起点。\n\n采访中，她谈到失败记录的重要性，也鼓励同学们保留对日常环境的敏感。",
			Author:       "沈言",
			Status:       StatusReturned,
			Edition:      "第 128 期",
			UpdatedLabel: "8 月 14 日",
		},
		{
			ID:           "interview-coach",
			Section:      SectionInterview,
			Title:        "教练席上的第十五个赛季",
			Summary:      "排球队主教练回顾队伍的成长与变化。",
			Body:         "清晨六点半，体育馆的灯准时亮起。对排球队主教练韩松来说，这是他在校队度过的第十五个赛季。",
			Author:       "许一帆",
			Status:       StatusPublished,
			Edition:      "第 127 期",
			UpdatedLabel: "8 月 12 日",
		},
		{
			ID:           "editorial-focus",
			Section:      SectionEditorial,
			Title:        "把课间十分钟还给专注",
			Summary:      "从短暂休息开始，重新理解学习节奏。",
			Body:         "课间并不是下一节课的预备铃。真正有效的学习，需要张弛有度的节奏，也需要短暂离开屏幕和书本的空间。",
			Author:       "编辑部",
			Status:       StatusDraft,
			Edition:      "第 128 期",
			UpdatedLabel: "今天 09:10",
		},
		{
			ID:           "event-music",
			Section:      SectionEvent,
			Title:        "秋季草坪音乐会节目单公布",
			Summary:      "十二组校园乐队与合唱团将在周五登台。",
			Body:         "秋季草坪音乐会将于本周五傍晚在东操场举行。节目单包含民乐合奏、无伴奏合唱和校园乐队演出。",
			Author:       "顾清禾",
			Status:       StatusPendingReview,
			Edition:      "第 128 期",
			UpdatedLabel: "今天 11:05",
		},
		{
			ID:           "archive-anniversary",
			Section:      SectionEvent,
			Title:        "百廿校庆专题报道",
			Summary:      "校庆系列报道归档稿。",
			Body:         "百廿校庆专题报道汇集了校史展览、校友返校和纪念大会的现场记录。本稿已随第 120 期完成归档。",
			Author:       "校庆报道组",
			Status:       StatusArchived,
			Edition:      "第 120 期",
			UpdatedLabel: "6 月 20 日",
		},
		{
			ID:           "event-volunteer",
			Section:      SectionEvent,
			Title:        "迎新志愿服务复盘完成",
			Summary:      "迎新志愿服务稿件已完成全部编务流程。",
			Body:         "本年度迎新志愿服务覆盖四个报到点，学生志愿者累计提供咨询、引导和行李转运服务。编务复盘现已完成。",
			Author:       "宋知遥",
			Status:       StatusCompleted,
			Edition:      "第 126 期",
			UpdatedLabel: "8 月 1 日",
		},
	}

	store := &Store{
		articles: make(map[string]Article, len(fixtures)),
		order:    make([]string, 0, len(fixtures)),
	}
	for _, article := range fixtures {
		store.articles[article.ID] = article
		store.order = append(store.order, article.ID)
	}
	return store
}

func (s *Store) List(section Section, status Status) []Article {
	s.mu.RLock()
	defer s.mu.RUnlock()

	articles := make([]Article, 0, len(s.articles))
	for _, id := range s.order {
		article := s.articles[id]
		if section != "" && article.Section != section {
			continue
		}
		if status != "" && article.Status != status {
			continue
		}
		articles = append(articles, article)
	}
	return articles
}

func (s *Store) Get(id string) (Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	article, ok := s.articles[id]
	if !ok {
		return Article{}, ErrNotFound
	}
	return article, nil
}

func (s *Store) UpdateContent(id string, update ContentUpdate) (Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	article, ok := s.articles[id]
	if !ok {
		return Article{}, ErrNotFound
	}
	article.Title = update.Title
	article.Summary = update.Summary
	article.Body = update.Body
	article.UpdatedLabel = "刚刚"
	s.articles[id] = article
	return article, nil
}

func (s *Store) SetStatus(id string, status Status) (Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	article, ok := s.articles[id]
	if !ok {
		return Article{}, ErrNotFound
	}
	article.Status = status
	article.UpdatedLabel = "刚刚"
	s.articles[id] = article
	return article, nil
}
