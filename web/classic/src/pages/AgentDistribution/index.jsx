/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Card, Col, Row, Table, Typography } from '@douyinfe/semi-ui';

import { API, showError, showSuccess } from '../../helpers';

const { Text, Title } = Typography;

const formatMoney = (value) => `¥${(Number(value || 0) / 100).toFixed(2)}`;

export default function AgentDistributionPage() {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [dashboard, setDashboard] = useState(null);
  const [records, setRecords] = useState([]);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [dashboardRes, recordsRes] = await Promise.all([
        API.get('/api/user/agent/dashboard?p=1&page_size=10'),
        API.get('/api/user/agent/records?p=1&page_size=10'),
      ]);
      if (!dashboardRes.data?.success) {
        throw new Error(
          dashboardRes.data?.message || t('获取代理分销信息失败'),
        );
      }
      setDashboard(dashboardRes.data.data);
      setRecords(recordsRes.data?.data?.items || []);
    } catch (error) {
      showError(error.message || t('获取代理分销信息失败'));
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const copyLink = async () => {
    if (!dashboard?.aff_code) return;
    await navigator.clipboard.writeText(
      `${window.location.origin}/register?aff=${dashboard.aff_code}`,
    );
    showSuccess(t('复制成功'));
  };

  const columns = [
    { title: t('来源'), dataIndex: 'source_type' },
    { title: t('用户'), dataIndex: 'invitee_name' },
    {
      title: t('分成比例'),
      dataIndex: 'commission_rate',
      render: (value) => `${Number(value || 0).toFixed(2)}%`,
    },
    {
      title: t('佣金'),
      dataIndex: 'commission_amount',
      render: formatMoney,
    },
  ];

  return (
    <div className='mt-[60px] px-2'>
      <div className='w-full max-w-7xl mx-auto space-y-4'>
        <div>
          <Title heading={3}>{t('代理分销')}</Title>
          <Text type='tertiary'>{t('邀请用户充值并获得首充与复充佣金')}</Text>
        </div>
        <Row gutter={16}>
          <Col span={8}>
            <Card title={t('可提现佣金')} loading={loading}>
              <Title heading={4}>
                {formatMoney(dashboard?.agent_commission_balance)}
              </Title>
            </Card>
          </Col>
          <Col span={8}>
            <Card title={t('累计佣金')} loading={loading}>
              <Title heading={4}>
                {formatMoney(dashboard?.agent_commission_total)}
              </Title>
            </Card>
          </Col>
          <Col span={8}>
            <Card title={t('代理用户')} loading={loading}>
              <Title heading={4}>{dashboard?.stats?.invitee_total || 0}</Title>
            </Card>
          </Col>
        </Row>
        <Card title={t('邀请链接')} loading={loading}>
          <div className='flex items-center gap-3'>
            <Text ellipsis={{ showTooltip: true }} className='flex-1'>
              {dashboard?.aff_code
                ? `${window.location.origin}/register?aff=${dashboard.aff_code}`
                : '-'}
            </Text>
            <Button theme='solid' onClick={copyLink}>
              {t('复制')}
            </Button>
          </div>
          <div className='mt-3'>
            <Text type='tertiary'>
              {t('首充分成比例')}：{dashboard?.effective_first_topup_rate || 0}
              %　
              {t('复充分成比例')}：{dashboard?.effective_repeat_topup_rate || 0}
              %
            </Text>
          </div>
        </Card>
        <Card title={t('佣金记录')}>
          <Table
            columns={columns}
            dataSource={records}
            loading={loading}
            pagination={false}
            rowKey='id'
          />
        </Card>
      </div>
    </div>
  );
}
