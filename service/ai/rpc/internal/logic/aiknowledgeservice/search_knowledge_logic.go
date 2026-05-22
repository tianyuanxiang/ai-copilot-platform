package aiknowledgeservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/model"
	"context"
	"fmt"
	"go-zero-rpc/common/xerr"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchKnowledgeLogic {
	return &SearchKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 执行知识库检索，支持 quick/deep 模式和不同检索范围。

func (l *SearchKnowledgeLogic) SearchKnowledge(in *pb.SearchKnowledgeReq) (*pb.SearchKnowledgeResp, error) {
	// 请求里传了 kbId
	// -> 查 ai_knowledge_base
	if in.HasKbId {
		kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
			}
			return nil, err
		}
		// -> 判断当前 userId 有没有读权限
		if !canAccessKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId) {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有访问该知识库的权限")
		}
		// 对用户问题做 query embedding
		querybedded, err := l.svcCtx.EngineCallClient.EngineEmbed(l.ctx, []string{in.Query})
		if err != nil {
			l.Logger.Errorf(
				"engine embed failed. user_id=%d kb_id=%d has_kb_id=%v domain_id=%d has_domain_id=%v err=%v",
				in.UserId,
				in.KbId,
				in.HasKbId,
				in.DomainId,
				in.HasDomainId,
				err,
			)
			return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "向量化服务暂时不可用，请稍后重试")
		}
		
		if len(querybedded.Vectors) != 1 {
			err = fmt.Errorf("embedding vector count mismatch: got %d, want %d", len(querybedded.Vectors), 1)
			return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "向量化服务返回的数量错误")
		}
	}

	// 请求里没传 kbId
	// -> 查出当前 userId 可访问的所有 kbId
	// -> accessibleKbIds = 查询结果
	// -> 如果为空：直接返回空 chunks
	return &pb.SearchKnowledgeResp{}, nil
}
