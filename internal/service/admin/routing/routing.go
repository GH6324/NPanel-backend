package routing

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/npanel-dev/NPanel-backend/api/admin/routing/v1"
	routingbiz "github.com/npanel-dev/NPanel-backend/internal/biz/admin/routing"
	publicrouting "github.com/npanel-dev/NPanel-backend/internal/biz/public/routing"
)

const successCode int32 = 200

type RoutingService struct {
	v1.UnimplementedRoutingServiceServer

	uc     *routingbiz.RoutingUsecase
	logger *log.Helper
}

func NewRoutingService(uc *routingbiz.RoutingUsecase, logger log.Logger) *RoutingService {
	return &RoutingService{uc: uc, logger: log.NewHelper(logger)}
}

func (s *RoutingService) BuildPublicConfig(ctx context.Context, now time.Time, opts publicrouting.ConfigOptions) (publicrouting.Envelope, error) {
	return s.uc.BuildConfig(ctx, now, opts)
}

func (s *RoutingService) RecordHealthReport(ctx context.Context, req publicrouting.HealthReportRequest) error {
	return s.uc.RecordHealthReport(ctx, req)
}

func (s *RoutingService) RecordRouteEvent(ctx context.Context, req publicrouting.RouteEventRequest) error {
	return s.uc.RecordRouteEvent(ctx, req)
}

func (s *RoutingService) ListRouteProfiles(ctx context.Context, req *v1.ListRouteProfilesRequest) (*v1.ListRouteProfilesReply, error) {
	items, total, err := s.uc.ListProfiles(ctx, int(req.Page), int(req.Size), req.Search, nil)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RouteProfile, 0, len(items))
	for _, item := range items {
		list = append(list, routeProfileToProto(item))
	}
	return &v1.ListRouteProfilesReply{Code: successCode, Message: "success", Data: &v1.RouteProfileListData{List: list, Total: total}}, nil
}

func (s *RoutingService) CreateRouteProfile(ctx context.Context, req *v1.CreateRouteProfileRequest) (*v1.RouteProfileReply, error) {
	item, err := s.uc.CreateProfile(ctx, routeProfileFromProto(req.GetProfile()))
	if err != nil {
		return nil, err
	}
	return &v1.RouteProfileReply{Code: successCode, Message: "success", Data: &v1.RouteProfileData{Profile: routeProfileToProto(item)}}, nil
}

func (s *RoutingService) UpdateRouteProfile(ctx context.Context, req *v1.UpdateRouteProfileRequest) (*v1.RouteProfileReply, error) {
	item, err := s.uc.UpdateProfile(ctx, routeProfileFromProto(req.GetProfile()))
	if err != nil {
		return nil, err
	}
	return &v1.RouteProfileReply{Code: successCode, Message: "success", Data: &v1.RouteProfileData{Profile: routeProfileToProto(item)}}, nil
}

func (s *RoutingService) DeleteRouteProfile(ctx context.Context, req *v1.DeleteRouteProfileRequest) (*v1.DeleteRouteItemReply, error) {
	if err := s.uc.DeleteProfile(ctx, req.Id); err != nil {
		return nil, err
	}
	return deleteReply(), nil
}

func (s *RoutingService) ListRouteRules(ctx context.Context, req *v1.ListRouteRulesRequest) (*v1.ListRouteRulesReply, error) {
	items, total, err := s.uc.ListRules(ctx, int(req.Page), int(req.Size), req.ProfileId, req.Search, nil)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RouteRule, 0, len(items))
	for _, item := range items {
		list = append(list, routeRuleToProto(item))
	}
	return &v1.ListRouteRulesReply{Code: successCode, Message: "success", Data: &v1.RouteRuleListData{List: list, Total: total}}, nil
}

func (s *RoutingService) CreateRouteRule(ctx context.Context, req *v1.CreateRouteRuleRequest) (*v1.RouteRuleReply, error) {
	item, err := s.uc.CreateRule(ctx, routeRuleFromProto(req.GetRule()))
	if err != nil {
		return nil, err
	}
	return &v1.RouteRuleReply{Code: successCode, Message: "success", Data: &v1.RouteRuleData{Rule: routeRuleToProto(item)}}, nil
}

func (s *RoutingService) UpdateRouteRule(ctx context.Context, req *v1.UpdateRouteRuleRequest) (*v1.RouteRuleReply, error) {
	item, err := s.uc.UpdateRule(ctx, routeRuleFromProto(req.GetRule()))
	if err != nil {
		return nil, err
	}
	return &v1.RouteRuleReply{Code: successCode, Message: "success", Data: &v1.RouteRuleData{Rule: routeRuleToProto(item)}}, nil
}

func (s *RoutingService) DeleteRouteRule(ctx context.Context, req *v1.DeleteRouteRuleRequest) (*v1.DeleteRouteItemReply, error) {
	if err := s.uc.DeleteRule(ctx, req.Id); err != nil {
		return nil, err
	}
	return deleteReply(), nil
}

func (s *RoutingService) ListDnsResolvers(ctx context.Context, req *v1.ListDnsResolversRequest) (*v1.ListDnsResolversReply, error) {
	items, total, err := s.uc.ListDNSResolvers(ctx, int(req.Page), int(req.Size), req.Search, nil)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.DnsResolver, 0, len(items))
	for _, item := range items {
		list = append(list, dnsResolverToProto(item))
	}
	return &v1.ListDnsResolversReply{Code: successCode, Message: "success", Data: &v1.DnsResolverListData{List: list, Total: total}}, nil
}

func (s *RoutingService) CreateDnsResolver(ctx context.Context, req *v1.CreateDnsResolverRequest) (*v1.DnsResolverReply, error) {
	item, err := s.uc.CreateDNSResolver(ctx, dnsResolverFromProto(req.GetResolver()))
	if err != nil {
		return nil, err
	}
	return &v1.DnsResolverReply{Code: successCode, Message: "success", Data: &v1.DnsResolverData{Resolver: dnsResolverToProto(item)}}, nil
}

func (s *RoutingService) UpdateDnsResolver(ctx context.Context, req *v1.UpdateDnsResolverRequest) (*v1.DnsResolverReply, error) {
	item, err := s.uc.UpdateDNSResolver(ctx, dnsResolverFromProto(req.GetResolver()))
	if err != nil {
		return nil, err
	}
	return &v1.DnsResolverReply{Code: successCode, Message: "success", Data: &v1.DnsResolverData{Resolver: dnsResolverToProto(item)}}, nil
}

func (s *RoutingService) DeleteDnsResolver(ctx context.Context, req *v1.DeleteDnsResolverRequest) (*v1.DeleteRouteItemReply, error) {
	if err := s.uc.DeleteDNSResolver(ctx, req.Id); err != nil {
		return nil, err
	}
	return deleteReply(), nil
}

func (s *RoutingService) ListRouteOutbounds(ctx context.Context, req *v1.ListRouteOutboundsRequest) (*v1.ListRouteOutboundsReply, error) {
	items, total, err := s.uc.ListOutbounds(ctx, int(req.Page), int(req.Size), req.Search, nil)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RouteOutbound, 0, len(items))
	for _, item := range items {
		list = append(list, routeOutboundToProto(item))
	}
	return &v1.ListRouteOutboundsReply{Code: successCode, Message: "success", Data: &v1.RouteOutboundListData{List: list, Total: total}}, nil
}

func (s *RoutingService) CreateRouteOutbound(ctx context.Context, req *v1.CreateRouteOutboundRequest) (*v1.RouteOutboundReply, error) {
	item, err := s.uc.CreateOutbound(ctx, routeOutboundFromProto(req.GetOutbound()))
	if err != nil {
		return nil, err
	}
	return &v1.RouteOutboundReply{Code: successCode, Message: "success", Data: &v1.RouteOutboundData{Outbound: routeOutboundToProto(item)}}, nil
}

func (s *RoutingService) UpdateRouteOutbound(ctx context.Context, req *v1.UpdateRouteOutboundRequest) (*v1.RouteOutboundReply, error) {
	item, err := s.uc.UpdateOutbound(ctx, routeOutboundFromProto(req.GetOutbound()))
	if err != nil {
		return nil, err
	}
	return &v1.RouteOutboundReply{Code: successCode, Message: "success", Data: &v1.RouteOutboundData{Outbound: routeOutboundToProto(item)}}, nil
}

func (s *RoutingService) DeleteRouteOutbound(ctx context.Context, req *v1.DeleteRouteOutboundRequest) (*v1.DeleteRouteItemReply, error) {
	if err := s.uc.DeleteOutbound(ctx, req.Id); err != nil {
		return nil, err
	}
	return deleteReply(), nil
}

func (s *RoutingService) ListUnlockServices(ctx context.Context, req *v1.ListUnlockServicesRequest) (*v1.ListUnlockServicesReply, error) {
	items, total, err := s.uc.ListUnlockServices(ctx, int(req.Page), int(req.Size), req.Search, nil)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.UnlockService, 0, len(items))
	for _, item := range items {
		list = append(list, unlockServiceToProto(item))
	}
	return &v1.ListUnlockServicesReply{Code: successCode, Message: "success", Data: &v1.UnlockServiceListData{List: list, Total: total}}, nil
}

func (s *RoutingService) CreateUnlockService(ctx context.Context, req *v1.CreateUnlockServiceRequest) (*v1.UnlockServiceReply, error) {
	item, err := s.uc.CreateUnlockService(ctx, unlockServiceFromProto(req.GetService()))
	if err != nil {
		return nil, err
	}
	return &v1.UnlockServiceReply{Code: successCode, Message: "success", Data: &v1.UnlockServiceData{Service: unlockServiceToProto(item)}}, nil
}

func (s *RoutingService) UpdateUnlockService(ctx context.Context, req *v1.UpdateUnlockServiceRequest) (*v1.UnlockServiceReply, error) {
	item, err := s.uc.UpdateUnlockService(ctx, unlockServiceFromProto(req.GetService()))
	if err != nil {
		return nil, err
	}
	return &v1.UnlockServiceReply{Code: successCode, Message: "success", Data: &v1.UnlockServiceData{Service: unlockServiceToProto(item)}}, nil
}

func (s *RoutingService) DeleteUnlockService(ctx context.Context, req *v1.DeleteUnlockServiceRequest) (*v1.DeleteRouteItemReply, error) {
	if err := s.uc.DeleteUnlockService(ctx, req.Id); err != nil {
		return nil, err
	}
	return deleteReply(), nil
}

func (s *RoutingService) PreviewRouteConfig(ctx context.Context, req *v1.PreviewRouteConfigRequest) (*v1.PreviewRouteConfigReply, error) {
	result, err := s.uc.Preview(ctx, routingbiz.PreviewRequest{
		Domain:            req.Domain,
		IP:                req.Ip,
		Port:              int(req.Port),
		UserID:            int64ToString(req.UserId),
		SubscribeID:       int64ToString(req.SubscribeId),
		UserSubscribeID:   int64ToString(req.UserSubscribeId),
		SubscribeToken:    req.SubscribeToken,
		NodeID:            int64ToString(req.NodeId),
		SupportedFeatures: req.SupportedFeatures,
	})
	if err != nil {
		return nil, err
	}
	ruleID, ruleName := "", ""
	if result.Rule != nil {
		ruleID = result.Rule.ID
		ruleName = result.Rule.Name
	}
	return &v1.PreviewRouteConfigReply{
		Code:    successCode,
		Message: "success",
		Data: &v1.PreviewRouteResult{
			RoutingHash:         result.RoutingHash,
			ProfileCode:         result.Profile.Code,
			ProfileName:         result.Profile.Name,
			ScopeType:           result.ScopeType,
			ScopeId:             result.ScopeID,
			Matched:             result.Matched,
			RuleId:              ruleID,
			RuleName:            ruleName,
			ActionType:          result.Action.Type,
			DnsResolverTag:      result.DNSResolverTag,
			OutboundTag:         result.OutboundTag,
			FallbackPolicy:      result.FallbackPolicy,
			UnsupportedFeatures: result.Unsupported,
			EffectiveMode:       result.EffectiveMode,
			ExecutionEnabled:    result.ExecutionEnabled,
		},
	}, nil
}

func (s *RoutingService) GetRoutingOverview(ctx context.Context, req *v1.GetRoutingOverviewRequest) (*v1.GetRoutingOverviewReply, error) {
	overview, err := s.uc.Overview(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoutingOverviewReply{
		Code:    successCode,
		Message: "success",
		Data:    routingOverviewToProto(overview),
	}, nil
}

func (s *RoutingService) ListRoutingHealthReports(ctx context.Context, req *v1.ListRoutingHealthReportsRequest) (*v1.ListRoutingHealthReportsReply, error) {
	items, total, err := s.uc.ListHealthReports(ctx, int(req.Page), int(req.Size), req.SubjectType, req.SubjectKey, req.ReporterType)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RoutingHealthReport, 0, len(items))
	for _, item := range items {
		list = append(list, healthReportToProto(item))
	}
	return &v1.ListRoutingHealthReportsReply{
		Code:    successCode,
		Message: "success",
		Data:    &v1.RoutingHealthReportListData{List: list, Total: total},
	}, nil
}

func (s *RoutingService) ListRoutingRouteEvents(ctx context.Context, req *v1.ListRoutingRouteEventsRequest) (*v1.ListRoutingRouteEventsReply, error) {
	items, total, err := s.uc.ListRouteEvents(ctx, int(req.Page), int(req.Size), req.EventType, req.ProfileCode, req.ReporterType)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RoutingRouteEvent, 0, len(items))
	for _, item := range items {
		list = append(list, routeEventToProto(item))
	}
	return &v1.ListRoutingRouteEventsReply{
		Code:    successCode,
		Message: "success",
		Data:    &v1.RoutingRouteEventListData{List: list, Total: total},
	}, nil
}

func (s *RoutingService) GetRoutingAnalytics(ctx context.Context, req *v1.GetRoutingAnalyticsRequest) (*v1.GetRoutingAnalyticsReply, error) {
	analytics, err := s.uc.Analytics(ctx, req.ProfileCode, req.RoutingHash, int(req.WindowMinutes))
	if err != nil {
		return nil, err
	}
	return &v1.GetRoutingAnalyticsReply{
		Code:    successCode,
		Message: "success",
		Data:    routingAnalyticsToProto(analytics),
	}, nil
}

func (s *RoutingService) ListRoutingGrayReleases(ctx context.Context, req *v1.ListRoutingGrayReleasesRequest) (*v1.ListRoutingGrayReleasesReply, error) {
	items, total, err := s.uc.ListGrayReleases(ctx, int(req.Page), int(req.Size), req.ProfileCode, req.Status)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.RoutingGrayRelease, 0, len(items))
	for _, item := range items {
		list = append(list, grayReleaseToProto(item))
	}
	return &v1.ListRoutingGrayReleasesReply{
		Code:    successCode,
		Message: "success",
		Data:    &v1.RoutingGrayReleaseListData{List: list, Total: total},
	}, nil
}

func (s *RoutingService) CreateRoutingGrayRelease(ctx context.Context, req *v1.CreateRoutingGrayReleaseRequest) (*v1.RoutingGrayReleaseReply, error) {
	item, err := s.uc.CreateGrayRelease(ctx, grayReleaseFromProto(req.GetRelease()))
	if err != nil {
		return nil, err
	}
	return &v1.RoutingGrayReleaseReply{Code: successCode, Message: "success", Data: &v1.RoutingGrayReleaseData{Release: grayReleaseToProto(item)}}, nil
}

func (s *RoutingService) UpdateRoutingGrayRelease(ctx context.Context, req *v1.UpdateRoutingGrayReleaseRequest) (*v1.RoutingGrayReleaseReply, error) {
	item, err := s.uc.UpdateGrayRelease(ctx, grayReleaseFromProto(req.GetRelease()))
	if err != nil {
		return nil, err
	}
	return &v1.RoutingGrayReleaseReply{Code: successCode, Message: "success", Data: &v1.RoutingGrayReleaseData{Release: grayReleaseToProto(item)}}, nil
}

func (s *RoutingService) DeleteRoutingGrayRelease(ctx context.Context, req *v1.DeleteRoutingGrayReleaseRequest) (*v1.DeleteRouteItemReply, error) {
	if err := s.uc.DeleteGrayRelease(ctx, req.Id); err != nil {
		return nil, err
	}
	return deleteReply(), nil
}

func (s *RoutingService) ActRoutingGrayRelease(ctx context.Context, req *v1.ActRoutingGrayReleaseRequest) (*v1.RoutingGrayReleaseReply, error) {
	item, err := s.uc.ActGrayRelease(ctx, req.Id, req.Action, req.Operator, req.Reason)
	if err != nil {
		return nil, err
	}
	return &v1.RoutingGrayReleaseReply{Code: successCode, Message: "success", Data: &v1.RoutingGrayReleaseData{Release: grayReleaseToProto(item)}}, nil
}

func routeProfileFromProto(item *v1.RouteProfile) *routingbiz.RouteProfile {
	if item == nil {
		return &routingbiz.RouteProfile{}
	}
	return &routingbiz.RouteProfile{
		ID:          item.Id,
		Code:        item.Code,
		Name:        item.Name,
		Description: item.Description,
		ScopeType:   item.ScopeType,
		ScopeID:     item.ScopeId,
		Priority:    int(item.Priority),
		Mode:        item.Mode,
		Enabled:     item.Enabled,
		ProfileJSON: item.ProfileJson,
	}
}

func healthReportToProto(item *routingbiz.RoutingHealthReport) *v1.RoutingHealthReport {
	if item == nil {
		return nil
	}
	return &v1.RoutingHealthReport{
		Id:                  item.ID,
		ReporterType:        item.ReporterType,
		ReporterId:          item.ReporterID,
		ProfileCode:         item.ProfileCode,
		RoutingHash:         item.RoutingHash,
		SubjectType:         item.SubjectType,
		SubjectKey:          item.SubjectKey,
		Region:              item.Region,
		Status:              item.Status,
		Source:              item.Source,
		RttMs:               int32(item.RTTMS),
		ConsecutiveFailures: int32(item.ConsecutiveFailures),
		LastError:           item.LastError,
		OutboundTag:         item.OutboundTag,
		DnsResolverTag:      item.DNSResolverTag,
		CheckedAt:           item.CheckedAt.Unix(),
		ReportJson:          item.ReportJSON,
		CreatedAt:           item.CreatedAt.Unix(),
		UpdatedAt:           item.UpdatedAt.Unix(),
	}
}

func routeEventToProto(item *routingbiz.RoutingRouteEvent) *v1.RoutingRouteEvent {
	if item == nil {
		return nil
	}
	return &v1.RoutingRouteEvent{
		Id:             item.ID,
		ReporterType:   item.ReporterType,
		ReporterId:     item.ReporterID,
		ProfileCode:    item.ProfileCode,
		RoutingHash:    item.RoutingHash,
		EventType:      item.EventType,
		Subject:        item.Subject,
		RuleId:         item.RuleID,
		RuleName:       item.RuleName,
		ActionType:     item.ActionType,
		OutboundTag:    item.OutboundTag,
		DnsResolverTag: item.DNSResolverTag,
		FallbackTarget: item.FallbackTarget,
		Status:         item.Status,
		LatencyMs:      int32(item.LatencyMS),
		Error:          item.Error,
		EventAt:        item.EventAt.Unix(),
		EventJson:      item.EventJSON,
		CreatedAt:      item.CreatedAt.Unix(),
		UpdatedAt:      item.UpdatedAt.Unix(),
	}
}

func grayReleaseFromProto(item *v1.RoutingGrayRelease) *routingbiz.RoutingGrayRelease {
	if item == nil {
		return &routingbiz.RoutingGrayRelease{}
	}
	return &routingbiz.RoutingGrayRelease{
		ID:             item.Id,
		ProfileCode:    item.ProfileCode,
		Name:           item.Name,
		Status:         item.Status,
		BatchNo:        int(item.BatchNo),
		TargetType:     item.TargetType,
		TargetIDsJSON:  item.TargetIdsJson,
		Operator:       item.Operator,
		RollbackReason: item.RollbackReason,
		StartedAt:      unixToTime(item.StartedAt),
		EndedAt:        unixToTime(item.EndedAt),
		ReleaseJSON:    item.ReleaseJson,
	}
}

func grayReleaseToProto(item *routingbiz.RoutingGrayRelease) *v1.RoutingGrayRelease {
	if item == nil {
		return nil
	}
	return &v1.RoutingGrayRelease{
		Id:             item.ID,
		ProfileCode:    item.ProfileCode,
		Name:           item.Name,
		Status:         item.Status,
		BatchNo:        int32(item.BatchNo),
		TargetType:     item.TargetType,
		TargetIdsJson:  item.TargetIDsJSON,
		Operator:       item.Operator,
		RollbackReason: item.RollbackReason,
		StartedAt:      unixOrZero(item.StartedAt),
		EndedAt:        unixOrZero(item.EndedAt),
		ReleaseJson:    item.ReleaseJSON,
		TargetCount:    int32(item.TargetCount),
		CreatedAt:      item.CreatedAt.Unix(),
		UpdatedAt:      item.UpdatedAt.Unix(),
	}
}

func routingAnalyticsToProto(item *routingbiz.RoutingAnalytics) *v1.RoutingAnalyticsData {
	if item == nil {
		return nil
	}
	items := make([]*v1.RoutingAnalyticsItem, 0, len(item.Items))
	for _, row := range item.Items {
		items = append(items, &v1.RoutingAnalyticsItem{
			ProfileCode:       row.ProfileCode,
			RoutingHash:       row.RoutingHash,
			ReporterId:        row.ReporterID,
			RouteEvents:       int32(row.RouteEvents),
			RouteDecisions:    int32(row.RouteDecisions),
			RouteFallbacks:    int32(row.RouteFallbacks),
			FallbackRateBp:    int32(row.FallbackRateBP),
			DnsFailures:       int32(row.DNSFailures),
			OutboundFailures:  int32(row.OutboundFailures),
			AffectedReporters: int32(row.AffectedReporters),
			LastEventType:     row.LastEventType,
			LastStatus:        row.LastStatus,
			LastError:         row.LastError,
			LastSeenAt:        unixOrZero(row.LastSeenAt),
		})
	}
	topErrors := make([]*v1.RoutingAnalyticsError, 0, len(item.TopErrors))
	for _, row := range item.TopErrors {
		topErrors = append(topErrors, &v1.RoutingAnalyticsError{
			Key:   row.Key,
			Kind:  row.Kind,
			Error: row.Error,
			Count: int32(row.Count),
		})
	}
	return &v1.RoutingAnalyticsData{
		Items:              items,
		TopErrors:          topErrors,
		TotalRouteEvents:   int32(item.TotalRouteEvents),
		TotalHealthReports: int32(item.TotalHealthReports),
		AffectedReporters:  int32(item.AffectedReporters),
		FallbackRateBp:     int32(item.FallbackRateBP),
		DnsFailRateBp:      int32(item.DNSFailRateBP),
		OutboundFailRateBp: int32(item.OutboundFailRateBP),
		WindowStartedAt:    unixOrZero(item.WindowStartedAt),
	}
}

func routeProfileToProto(item *routingbiz.RouteProfile) *v1.RouteProfile {
	if item == nil {
		return nil
	}
	return &v1.RouteProfile{
		Id:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		Description: item.Description,
		ScopeType:   item.ScopeType,
		ScopeId:     item.ScopeID,
		Priority:    int32(item.Priority),
		Mode:        item.Mode,
		Enabled:     item.Enabled,
		ProfileJson: item.ProfileJSON,
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}
}

func routeRuleFromProto(item *v1.RouteRule) *routingbiz.RouteRule {
	if item == nil {
		return &routingbiz.RouteRule{}
	}
	return &routingbiz.RouteRule{
		ID:          item.Id,
		ProfileID:   item.ProfileId,
		Name:        item.Name,
		Priority:    int(item.Priority),
		Enabled:     item.Enabled,
		ServiceCode: item.ServiceCode,
		MatcherJSON: item.MatcherJson,
		ActionJSON:  item.ActionJson,
	}
}

func routeRuleToProto(item *routingbiz.RouteRule) *v1.RouteRule {
	if item == nil {
		return nil
	}
	return &v1.RouteRule{
		Id:          item.ID,
		ProfileId:   item.ProfileID,
		Name:        item.Name,
		Priority:    int32(item.Priority),
		Enabled:     item.Enabled,
		ServiceCode: item.ServiceCode,
		MatcherJson: item.MatcherJSON,
		ActionJson:  item.ActionJSON,
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}
}

func dnsResolverFromProto(item *v1.DnsResolver) *routingbiz.DNSResolver {
	if item == nil {
		return &routingbiz.DNSResolver{}
	}
	return &routingbiz.DNSResolver{
		ID:           item.Id,
		Tag:          item.Tag,
		Name:         item.Name,
		Proto:        item.Proto,
		Address:      item.Address,
		Port:         int(item.Port),
		Enabled:      item.Enabled,
		ResolverJSON: item.ResolverJson,
	}
}

func dnsResolverToProto(item *routingbiz.DNSResolver) *v1.DnsResolver {
	if item == nil {
		return nil
	}
	return &v1.DnsResolver{
		Id:           item.ID,
		Tag:          item.Tag,
		Name:         item.Name,
		Proto:        item.Proto,
		Address:      item.Address,
		Port:         int32(item.Port),
		Enabled:      item.Enabled,
		ResolverJson: item.ResolverJSON,
		CreatedAt:    item.CreatedAt.Unix(),
		UpdatedAt:    item.UpdatedAt.Unix(),
	}
}

func routeOutboundFromProto(item *v1.RouteOutbound) *routingbiz.RouteOutbound {
	if item == nil {
		return &routingbiz.RouteOutbound{}
	}
	return &routingbiz.RouteOutbound{
		ID:           item.Id,
		Tag:          item.Tag,
		Name:         item.Name,
		Type:         item.Type,
		Region:       item.Region,
		Enabled:      item.Enabled,
		OutboundJSON: item.OutboundJson,
	}
}

func routeOutboundToProto(item *routingbiz.RouteOutbound) *v1.RouteOutbound {
	if item == nil {
		return nil
	}
	return &v1.RouteOutbound{
		Id:           item.ID,
		Tag:          item.Tag,
		Name:         item.Name,
		Type:         item.Type,
		Region:       item.Region,
		Enabled:      item.Enabled,
		OutboundJson: item.OutboundJSON,
		CreatedAt:    item.CreatedAt.Unix(),
		UpdatedAt:    item.UpdatedAt.Unix(),
	}
}

func unlockServiceFromProto(item *v1.UnlockService) *routingbiz.UnlockService {
	if item == nil {
		return &routingbiz.UnlockService{}
	}
	return &routingbiz.UnlockService{
		ID:          item.Id,
		Code:        item.Code,
		Name:        item.Name,
		Category:    item.Category,
		Enabled:     item.Enabled,
		ServiceJSON: item.ServiceJson,
	}
}

func unlockServiceToProto(item *routingbiz.UnlockService) *v1.UnlockService {
	if item == nil {
		return nil
	}
	return &v1.UnlockService{
		Id:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		Category:    item.Category,
		Enabled:     item.Enabled,
		ServiceJson: item.ServiceJSON,
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}
}

func deleteReply() *v1.DeleteRouteItemReply {
	return &v1.DeleteRouteItemReply{Code: successCode, Message: "success", Data: &v1.DeleteRouteItemData{Success: true}}
}

func routingOverviewToProto(item *routingbiz.RoutingOverview) *v1.RoutingOverview {
	if item == nil {
		return nil
	}
	health := make([]*v1.RoutingHealthItem, 0, len(item.Health))
	for _, healthItem := range item.Health {
		health = append(health, &v1.RoutingHealthItem{
			Kind:                healthItem.Kind,
			Key:                 healthItem.Key,
			Name:                healthItem.Name,
			Status:              healthItem.Status,
			Source:              healthItem.Source,
			CheckedAt:           healthItem.CheckedAt.Unix(),
			RttMs:               int32(healthItem.RTTMS),
			ConsecutiveFailures: int32(healthItem.ConsecutiveFailures),
			LastError:           healthItem.LastError,
			OutboundTag:         healthItem.OutboundTag,
			DnsResolverTag:      healthItem.DNSResolverTag,
		})
	}
	guards := make([]*v1.RoutingEnforceGuard, 0, len(item.Guards))
	for _, guard := range item.Guards {
		guards = append(guards, &v1.RoutingEnforceGuard{
			Key:    guard.Key,
			Label:  guard.Label,
			Passed: guard.Passed,
			Status: guard.Status,
			Reason: guard.Reason,
		})
	}
	auditEvents := make([]*v1.RoutingAuditEvent, 0, len(item.AuditEvents))
	for _, event := range item.AuditEvents {
		auditEvents = append(auditEvents, &v1.RoutingAuditEvent{
			Id:           event.ID,
			ResourceType: event.ResourceType,
			ResourceId:   event.ResourceID,
			ResourceName: event.ResourceName,
			Action:       event.Action,
			Summary:      event.Summary,
			CreatedAt:    event.CreatedAt.Unix(),
		})
	}
	return &v1.RoutingOverview{
		RoutingHash:      item.RoutingHash,
		GeneratedAt:      item.GeneratedAt,
		ProfileCode:      item.ProfileCode,
		ProfileName:      item.ProfileName,
		Mode:             item.Mode,
		ProfileEnabled:   item.ProfileEnabled,
		EnforceReady:     item.EnforceReady,
		ExecutionEnabled: item.ExecutionEnabled,
		RollbackAction:   item.RollbackAction,
		CompileError:     item.CompileError,
		Health:           health,
		Guards:           guards,
		AuditEvents:      auditEvents,
	}
}

func int64ToString(value int64) string {
	if value == 0 {
		return ""
	}
	return strconv.FormatInt(value, 10)
}

func unixToTime(value int64) time.Time {
	if value <= 0 {
		return time.Time{}
	}
	return time.Unix(value, 0)
}

func unixOrZero(value time.Time) int64 {
	if value.IsZero() {
		return 0
	}
	return value.Unix()
}
