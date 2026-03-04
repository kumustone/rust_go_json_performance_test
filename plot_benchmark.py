#!/usr/bin/env python3
"""
WafHttpRequest JSON 序列化性能对比可视化
"""

import matplotlib.pyplot as plt
import numpy as np

# 数据来源：Go 和 Rust 基准测试结果
# 时间单位：微秒 (µs)

data = {
    'S': {
        'encode': {
            'go_easyjson': 1.201,
            'go_stdlib': 1.491,
            'go_jsonv2': 2.644,
            'rust_serde': 1.019,
        },
        'decode': {
            'go_easyjson': 1.565,
            'go_stdlib': 7.267,
            'go_jsonv2': 2.784,
            'rust_serde': 1.615,
        }
    },
    'M': {
        'encode': {
            'go_easyjson': 13.886,
            'go_stdlib': 7.968,
            'go_jsonv2': 23.285,
            'rust_serde': 7.896,
        },
        'decode': {
            'go_easyjson': 7.387,
            'go_stdlib': 58.173,
            'go_jsonv2': 15.844,
            'rust_serde': 7.012,
        }
    },
    'L': {
        'encode': {
            'go_easyjson': 97.555,
            'go_stdlib': 47.228,
            'go_jsonv2': 159.052,
            'rust_serde': 59.057,
        },
        'decode': {
            'go_easyjson': 43.138,
            'go_stdlib': 412.086,
            'go_jsonv2': 99.297,
            'rust_serde': 41.502,
        }
    }
}

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['Arial Unicode MS', 'SimHei', 'DejaVu Sans']
plt.rcParams['axes.unicode_minus'] = False

# 创建图表
fig, axes = plt.subplots(2, 3, figsize=(16, 10))
fig.suptitle('WafHttpRequest JSON 序列化性能对比 (时间越短越好)', fontsize=16, fontweight='bold')

cases = ['S', 'M', 'L']
operations = ['encode', 'decode']
labels = ['Go easyjson', 'Go stdlib', 'Go json/v2', 'Rust serde_json']
colors = ['#3498db', '#e74c3c', '#f39c12', '#2ecc71']

x = np.arange(len(labels))
width = 0.7

for i, case in enumerate(cases):
    for j, op in enumerate(operations):
        ax = axes[j, i]
        
        values = [
            data[case][op]['go_easyjson'],
            data[case][op]['go_stdlib'],
            data[case][op]['go_jsonv2'],
            data[case][op]['rust_serde']
        ]
        
        bars = ax.bar(x, values, width, color=colors, alpha=0.8, edgecolor='black', linewidth=1.2)
        
        # 添加数值标签
        for bar, val in zip(bars, values):
            height = bar.get_height()
            ax.text(bar.get_x() + bar.get_width()/2., height,
                   f'{val:.2f}',
                   ha='center', va='bottom', fontsize=9, fontweight='bold')
        
        # 设置标题和标签
        op_cn = '序列化' if op == 'encode' else '反序列化'
        case_desc = {
            'S': 'Small (512B)',
            'M': 'Medium (8KB)',
            'L': 'Large (64KB)'
        }
        ax.set_title(f'{case_desc[case]} - {op_cn}', fontsize=12, fontweight='bold')
        ax.set_ylabel('时间 (微秒)', fontsize=10)
        ax.set_xticks(x)
        ax.set_xticklabels(labels, rotation=15, ha='right', fontsize=9)
        ax.grid(axis='y', alpha=0.3, linestyle='--')
        
        # 对于 decode 的 L 档，stdlib 数值太大，单独标注
        if case == 'L' and op == 'decode':
            ax.set_ylim(0, max(values) * 1.2)

plt.tight_layout()
plt.savefig('benchmark_comparison.png', 
            dpi=300, bbox_inches='tight')
print("图表已保存到: benchmark_comparison.png")

# 创建第二个图表：性能倍率对比
fig2, axes2 = plt.subplots(1, 2, figsize=(14, 6))
fig2.suptitle('性能倍率对比 (相对于 Go stdlib，倍率越高越快)', fontsize=14, fontweight='bold')

for j, op in enumerate(operations):
    ax = axes2[j]
    
    # 计算相对于 stdlib 的倍率
    ratios = {
        'go_easyjson': [],
        'go_jsonv2': [],
        'rust_serde': []
    }
    
    for case in cases:
        stdlib_time = data[case][op]['go_stdlib']
        ratios['go_easyjson'].append(stdlib_time / data[case][op]['go_easyjson'])
        ratios['go_jsonv2'].append(stdlib_time / data[case][op]['go_jsonv2'])
        ratios['rust_serde'].append(stdlib_time / data[case][op]['rust_serde'])
    
    x_pos = np.arange(len(cases))
    width = 0.25
    
    bars1 = ax.bar(x_pos - width, ratios['go_easyjson'], width, 
                   label='Go easyjson', color='#3498db', alpha=0.8, edgecolor='black')
    bars2 = ax.bar(x_pos, ratios['go_jsonv2'], width,
                   label='Go json/v2', color='#f39c12', alpha=0.8, edgecolor='black')
    bars3 = ax.bar(x_pos + width, ratios['rust_serde'], width,
                   label='Rust serde_json', color='#2ecc71', alpha=0.8, edgecolor='black')
    
    # 添加数值标签
    for bars in [bars1, bars2, bars3]:
        for bar in bars:
            height = bar.get_height()
            ax.text(bar.get_x() + bar.get_width()/2., height,
                   f'{height:.2f}x',
                   ha='center', va='bottom', fontsize=8, fontweight='bold')
    
    ax.axhline(y=1.0, color='red', linestyle='--', linewidth=2, alpha=0.5, label='Go stdlib (基准)')
    
    op_cn = '序列化' if op == 'encode' else '反序列化'
    ax.set_title(f'{op_cn}性能倍率', fontsize=12, fontweight='bold')
    ax.set_ylabel('相对速度倍率', fontsize=10)
    ax.set_xlabel('数据规模', fontsize=10)
    ax.set_xticks(x_pos)
    ax.set_xticklabels(['S (512B)', 'M (8KB)', 'L (64KB)'])
    ax.legend(fontsize=9)
    ax.grid(axis='y', alpha=0.3, linestyle='--')

plt.tight_layout()
plt.savefig('benchmark_ratio.png',
            dpi=300, bbox_inches='tight')
print("倍率对比图已保存到: benchmark_ratio.png")

print("\n性能总结：")
print("=" * 60)
for case in cases:
    print(f"\n{case} 档位:")
    for op in operations:
        op_cn = '序列化' if op == 'encode' else '反序列化'
        values = data[case][op]
        fastest = min(values.items(), key=lambda x: x[1])
        print(f"  {op_cn}: 最快 = {fastest[0]} ({fastest[1]:.2f} µs)")
