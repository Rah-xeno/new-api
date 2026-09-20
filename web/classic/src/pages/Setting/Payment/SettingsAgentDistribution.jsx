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

import React, { useEffect, useRef, useState } from 'react';
import {
  Banner,
  Button,
  Col,
  Form,
  Row,
  Spin,
  Typography,
} from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { API, showError, showSuccess } from '../../../helpers';

const { Text } = Typography;
const FIRST_RATE_KEY = 'agent_distribution_setting.default_first_topup_rate';
const REPEAT_RATE_KEY = 'agent_distribution_setting.default_repeat_topup_rate';

export default function SettingsAgentDistribution({ options, refresh }) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({ firstRate: 0, repeatRate: 0 });
  const formApiRef = useRef(null);

  useEffect(() => {
    const values = {
      firstRate: Number(options?.[FIRST_RATE_KEY] ?? 0),
      repeatRate: Number(options?.[REPEAT_RATE_KEY] ?? 0),
    };
    setInputs(values);
    formApiRef.current?.setValues(values);
  }, [options]);

  const submitSettings = async () => {
    const firstRate = Number(inputs.firstRate || 0);
    const repeatRate = Number(inputs.repeatRate || 0);
    if (
      firstRate < 0 ||
      firstRate > 100 ||
      repeatRate < 0 ||
      repeatRate > 100
    ) {
      showError(t('代理分销比例必须在 0 到 100 之间'));
      return;
    }

    setLoading(true);
    try {
      const responses = await Promise.all([
        API.put('/api/option/', {
          key: FIRST_RATE_KEY,
          value: String(firstRate),
        }),
        API.put('/api/option/', {
          key: REPEAT_RATE_KEY,
          value: String(repeatRate),
        }),
      ]);
      const failed = responses.find((response) => !response.data.success);
      if (failed) {
        showError(failed.data.message || t('更新失败'));
        return;
      }
      showSuccess(t('更新成功'));
      await refresh?.();
    } catch {
      showError(t('更新失败'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Spin spinning={loading}>
      <Form
        initValues={inputs}
        getFormApi={(api) => (formApiRef.current = api)}
        onValueChange={setInputs}
      >
        <Form.Section text={t('代理分销默认配置')}>
          <Text>
            {t(
              '这里配置代理用户的默认首充和复充分成比例，单个代理也可以使用自己的比例。',
            )}
          </Text>
          <Banner
            type='info'
            description={t(
              '代理邀请的用户可以多次产生返利；首笔成功支付使用首充比例，之后在代理有效期间持续使用复充比例。',
            )}
          />
          <Row gutter={16}>
            <Col span={12}>
              <Form.InputNumber
                field='firstRate'
                label={t('默认首充分成比例')}
                min={0}
                max={100}
                precision={2}
                suffix='%'
                style={{ width: '100%' }}
              />
            </Col>
            <Col span={12}>
              <Form.InputNumber
                field='repeatRate'
                label={t('默认复充分成比例')}
                min={0}
                max={100}
                precision={2}
                suffix='%'
                style={{ width: '100%' }}
              />
            </Col>
          </Row>
          <Button onClick={submitSettings}>{t('保存代理分销默认配置')}</Button>
        </Form.Section>
      </Form>
    </Spin>
  );
}
