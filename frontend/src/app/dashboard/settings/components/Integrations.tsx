// frontend/src/app/dashboard/settings/components/Integrations.tsx
'use client';

import React, { useState, useEffect } from 'react';
import {
  Integration,
  IntegrationCreate,
  listIntegrations,
  createIntegration,
  deleteIntegration,
  testIntegration,
} from '@/lib/api/settings';
import {
  Plus,
  Trash2,
  CheckCircle,
  XCircle,
  AlertCircle,
  RefreshCw
} from 'lucide-react';

const INTEGRATION_PROVIDERS = [
  {
    id: 'google_ads',
    name: 'Google Ads',
    description: 'Import and sync campaigns from Google Ads',
    icon: '🎯',
    authType: 'oauth',
    fields: ['access_token', 'refresh_token'],
    category: 'advertising',
  },
  {
    id: 'facebook_ads',
    name: 'Facebook Ads',
    description: 'Manage Meta advertising campaigns',
    icon: '📘',
    authType: 'oauth',
    fields: ['access_token'],
    category: 'advertising',
  },
  {
    id: 'microsoft_ads',
    name: 'Microsoft Advertising',
    description: 'Connect with Bing Ads platform',
    icon: '🔷',
    authType: 'oauth',
    fields: ['access_token', 'refresh_token'],
    category: 'advertising',
  },
  {
    id: 'slack',
    name: 'Slack',
    description: 'Receive alerts in your Slack workspace',
    icon: '💬',
    authType: 'webhook',
    fields: ['webhook_url'],
    category: 'communication',
  },
  {
    id: 'webhook',
    name: 'Custom Webhook',
    description: 'Send events to your custom endpoint',
    icon: '🔗',
    authType: 'webhook',
    fields: ['webhook_url'],
    category: 'communication',
  },
  {
    id: 'google_analytics',
    name: 'Google Analytics',
    description: 'Enhanced tracking and attribution',
    icon: '📊',
    authType: 'api_key',
    fields: ['api_key'],
    category: 'analytics',
  },
  {
    id: 'segment',
    name: 'Segment',
    description: 'Customer data platform integration',
    icon: '🎨',
    authType: 'api_key',
    fields: ['api_key'],
    category: 'analytics',
  },
  {
    id: 'stripe',
    name: 'Stripe',
    description: 'Payment processing and billing',
    icon: '💳',
    authType: 'api_key',
    fields: ['api_key', 'api_secret'],
    category: 'billing',
  },
  {
    id: 'sendgrid',
    name: 'SendGrid',
    description: 'Email delivery service',
    icon: '📧',
    authType: 'api_key',
    fields: ['api_key'],
    category: 'communication',
  },
];

export function Integrations() {
  const [integrations, setIntegrations] = useState<Integration[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [testingId, setTestingId] = useState<string | null>(null);
  console.log('selectedProvider', selectedProvider);

  useEffect(() => {
    loadIntegrations();
  }, []);

  const loadIntegrations = async () => {
    try {
      const data = await listIntegrations();
      setIntegrations(data.integrations);
    } catch (error) {
      console.error('Failed to load integrations:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this integration?')) return;

    try {
      await deleteIntegration(id);
      await loadIntegrations();
    } catch (error) {
      console.error('Failed to delete integration:', error);
      alert('Failed to delete integration');
    }
  };

  const handleTest = async (id: string) => {
    setTestingId(id);
    try {
      await testIntegration(id);
      alert('Integration test successful!');
      await loadIntegrations();
    } catch (error) {
      console.error('Integration test failed:', error);
      alert('Integration test failed');
    } finally {
      setTestingId(null);
    }
  };

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      active: { icon: CheckCircle, color: 'text-green-600 bg-green-50', label: 'Active' },
      inactive: { icon: XCircle, color: 'text-gray-600 bg-gray-50', label: 'Inactive' },
      error: { icon: AlertCircle, color: 'text-red-600 bg-red-50', label: 'Error' },
      expired: { icon: AlertCircle, color: 'text-yellow-600 bg-yellow-50', label: 'Expired' },
    };

    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.inactive;
    const Icon = config.icon;

    return (
      <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${config.color}`}>
        <Icon className="w-3 h-3" />
        {config.label}
      </span>
    );
  };

  const getCategoryColor = (category: string) => {
    const colors = {
      advertising: 'bg-blue-100 text-blue-700',
      communication: 'bg-green-100 text-green-700',
      analytics: 'bg-purple-100 text-purple-700',
      billing: 'bg-orange-100 text-orange-700',
    };
    return colors[category as keyof typeof colors] || 'bg-gray-100 text-gray-700';
  };

  if (loading) {
    return <div className="text-center py-8">Loading integrations...</div>;
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <p className="text-gray-600 dark:text-slate-400">
          Connect third-party services to enhance your bid optimization platform
        </p>
        <button
          onClick={() => setShowAddModal(true)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-4 h-4" />
          Add Integration
        </button>
      </div>

      {/* Integrations List */}
      {integrations?.length === 0 ? (
        <div className="text-center py-12 bg-gray-50 dark:bg-slate-800/50 rounded-lg border-2 border-dashed border-gray-300 dark:border-slate-700">
          <div className="text-4xl mb-4">🔌</div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-slate-100 mb-2">
            No integrations yet
          </h3>
          <p className="text-gray-600 dark:text-slate-400 mb-4">
            Connect your first integration to get started
          </p>
          <button
            onClick={() => setShowAddModal(true)}
            className="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Add Your First Integration
          </button>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {integrations && integrations.length > 0 && integrations.map((integration) => {
            const provider = INTEGRATION_PROVIDERS.find((p) => p.id === integration.provider);
            return (
              <div
                key={integration.id}
                className="bg-white dark:bg-slate-800 rounded-lg border border-gray-200 dark:border-slate-700 p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="text-2xl">{provider?.icon || '🔌'}</div>
                    <div>
                      <h4 className="font-semibold text-gray-900 dark:text-slate-100">
                        {integration.integration_name}
                      </h4>
                      <p className="text-sm text-gray-600 dark:text-slate-400">
                        {provider?.name || integration.provider}
                      </p>
                    </div>
                  </div>
                  {getStatusBadge(integration.status)}
                </div>

                {provider && (
                  <span className={`inline-block px-2 py-1 text-xs rounded-full ${getCategoryColor(provider.category)}`}>
                    {provider.category}
                  </span>
                )}

                {integration.last_sync_at && (
                  <p className="text-xs text-gray-500 dark:text-slate-500 mt-2">
                    Last synced: {new Date(integration.last_sync_at).toLocaleString()}
                  </p>
                )}

                {integration.last_error && (
                  <p className="text-xs text-red-600 dark:text-red-400 mt-2 flex items-center gap-1">
                    <AlertCircle className="w-3 h-3" />
                    {integration.last_error}
                  </p>
                )}

                <div className="flex gap-2 mt-4 pt-4 border-t border-gray-200 dark:border-slate-700">
                  <button
                    onClick={() => handleTest(integration.id)}
                    disabled={testingId === integration.id}
                    className="flex-1 px-3 py-1.5 text-sm border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors disabled:opacity-50"
                  >
                    <RefreshCw className={`w-4 h-4 mx-auto ${testingId === integration.id ? 'animate-spin' : ''}`} />
                  </button>
                  <button
                    onClick={() => handleDelete(integration.id)}
                    className="flex-1 px-3 py-1.5 text-sm border border-red-300 text-red-600 dark:border-red-700 dark:text-red-400 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                  >
                    <Trash2 className="w-4 h-4 mx-auto" />
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Add Integration Modal */}
      {showAddModal && (
        <AddIntegrationModal
          providers={INTEGRATION_PROVIDERS}
          onClose={() => {
            setShowAddModal(false);
            setSelectedProvider(null);
          }}
          onSuccess={() => {
            setShowAddModal(false);
            setSelectedProvider(null);
            loadIntegrations();
          }}
        />
      )}
    </div>
  );
}

// Add Integration Modal Component
function AddIntegrationModal({
  providers,
  onClose,
  onSuccess,
}: {
  providers: typeof INTEGRATION_PROVIDERS;
  onClose: () => void;
  onSuccess: () => void;
}) {
  const [step, setStep] = useState<'select' | 'configure'>('select');
  const [selectedProvider, setSelectedProvider] = useState<typeof INTEGRATION_PROVIDERS[0] | null>(null);
  const [formData, setFormData] = useState<Record<string, string>>({
    integration_name: '',
  });
  const [saving, setSaving] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  const categories = ['all', ...new Set(providers.map((p) => p.category))];
  const filteredProviders =
    selectedCategory === 'all'
      ? providers
      : providers.filter((p) => p.category === selectedCategory);

  const handleProviderSelect = (provider: typeof INTEGRATION_PROVIDERS[0]) => {
    setSelectedProvider(provider);
    setFormData({ integration_name: provider.name });
    setStep('configure');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProvider) return;

    setSaving(true);
    try {
      const integrationData: IntegrationCreate = {
        provider: selectedProvider.id,
        integration_name: formData.integration_name,
        auth_type: selectedProvider.authType,
        ...(formData.access_token && { access_token: formData.access_token }),
        ...(formData.refresh_token && { refresh_token: formData.refresh_token }),
        ...(formData.api_key && { api_key: formData.api_key }),
        ...(formData.api_secret && { api_secret: formData.api_secret }),
        ...(formData.webhook_url && { webhook_url: formData.webhook_url }),
      };

      await createIntegration(integrationData);
      onSuccess();
    } catch (error) {
      console.error('Failed to create integration:', error);
      alert('Failed to create integration');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-slate-800 rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6 border-b border-gray-200 dark:border-slate-700">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-slate-100">
            {step === 'select' ? 'Select Integration' : `Configure ${selectedProvider?.name}`}
          </h3>
        </div>

        <div className="p-6">
          {step === 'select' ? (
            <div className="space-y-4">
              {/* Category Filter */}
              <div className="flex gap-2 flex-wrap">
                {categories.map((cat) => (
                  <button
                    key={cat}
                    onClick={() => setSelectedCategory(cat)}
                    className={`px-3 py-1 rounded-full text-sm capitalize ${
                      selectedCategory === cat
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-slate-300'
                    }`}
                  >
                    {cat}
                  </button>
                ))}
              </div>

              {/* Provider Grid */}
              <div className="grid gap-3 md:grid-cols-2">
                {filteredProviders.map((provider) => (
                  <button
                    key={provider.id}
                    onClick={() => handleProviderSelect(provider)}
                    className="text-left p-4 border border-gray-200 dark:border-slate-700 rounded-lg hover:border-blue-500 hover:shadow-md transition-all"
                  >
                    <div className="flex items-start gap-3">
                      <div className="text-2xl">{provider.icon}</div>
                      <div className="flex-1">
                        <h4 className="font-medium text-gray-900 dark:text-slate-100">
                          {provider.name}
                        </h4>
                        <p className="text-sm text-gray-600 dark:text-slate-400 mt-1">
                          {provider.description}
                        </p>
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
                  Integration Name
                </label>
                <input
                  type="text"
                  value={formData.integration_name}
                  onChange={(e) => setFormData({ ...formData, integration_name: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
                  required
                />
              </div>

              {selectedProvider?.fields.map((field) => (
                <div key={field}>
                  <label className="block text-sm font-medium text-gray-700 dark:text-slate-300 mb-2">
                    {field.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')}
                  </label>
                  <input
                    type={field.includes('secret') || field.includes('token') ? 'password' : 'text'}
                    value={formData[field] || ''}
                    onChange={(e) => setFormData({ ...formData, [field]: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-slate-700 dark:text-slate-100"
                    placeholder={field.includes('url') ? 'https://...' : ''}
                    required
                  />
                </div>
              ))}

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setStep('select')}
                  className="flex-1 px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
                >
                  Back
                </button>
                <button
                  type="submit"
                  disabled={saving}
                  className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
                >
                  {saving ? 'Creating...' : 'Create Integration'}
                </button>
              </div>
            </form>
          )}
        </div>

        <div className="p-4 border-t border-gray-200 dark:border-slate-700">
          <button
            onClick={onClose}
            className="w-full px-4 py-2 text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-100 transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}
