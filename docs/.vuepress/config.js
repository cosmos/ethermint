module.exports = {
  theme: 'cosmos',
  title: 'Ethermint Documentation',
  locales: {
    '/': {
      lang: 'en-US'
    },
  },
  base: process.env.VUEPRESS_BASE || '/',
  themeConfig: {
    repo: 'cosmos/ethermint',
    docsRepo: 'cosmos/ethermint',
    docsDir: 'docs',
    editLinks: true,
    custom: true,
    logo: {
      src: '/logo.svg',
    },
    algolia: {
      id: 'BH4D9OD16A',
      key: 'ac317234e6a42074175369b2f42e9754',
      index: 'ethermint'
    },
    sidebar: { 
      auto: true,
      nav: [
        {
          title: 'Introduction',
          children: [
            {
              title: 'High-Level Overview',
              path: '/intro/overview.html'
            },
            {
              title: 'Architecture',
              path: '/intro/architecture.html'
            }
          ]
        },
        {
          title: 'Basics',
          children: [
            {
              title: 'Accounts',
              path: '/basics/accounts.html'
            },
            {
              title: 'Transactions',
              path: '/basics/transactions.html'
            },
            {
              title: 'Gas',
              path: '/basics/gas.html'
            }
          ]
        },
        {
          title: 'Core Concepts',
          children: [
            {
              title: 'Encoding',
              path: '/core/encoding.html'
            },
            {
              title: 'Events',
              path: '/core/events.html'
            },
          ]
        },
        {
          title: 'Guides',
          children: [
            {
              title: 'Clients',
              path: '/clients'
            }
          ]
        },
        {
          title: 'Specifications',
          children: [
            {
              title: 'Modules',
              directory: true,
              path: '/modules'
            }
          ]
        },
        {
          title: 'Resources',
          children: [
            {
              title: 'Ethermint API Reference',
              path: 'https://godoc.org/github.com/cosmos/ethermint'
            },
            {
              title: 'Cosmos REST API Spec',
              path: 'https://cosmos.network/rpc/'
            },
            {
              title: 'Ethereum JSON RPC API Reference',
              path: 'https://eth.wiki/json-rpc/API'
            }
          ]
        }
      ]
    },
    gutter: {
      title: 'Help & Support',
      editLink: true,
      chat: {
        title: 'Developer Chat',
        text: 'Chat with Ethermint developers on Discord.',
        url: 'https://discordapp.com/channels/669268347736686612',
        bg: 'linear-gradient(103.75deg, #1B1E36 0%, #22253F 100%)'
      },
      forum: {
        title: 'Ethermint Developer Forum',
        text: 'Join the Ethermint Developer Forum to learn more.',
        url: 'https://forum.cosmos.network/',
        bg: 'linear-gradient(221.79deg, #3D6B99 -1.08%, #336699 95.88%)',
        logo: 'ethereum-white'
      },
      github: {
        title: 'Found an Issue?',
        text: 'Help us improve this page by suggesting edits on GitHub.',
        bg: '#F8F9FC'
      }
    },
    footer: {
      logo: '/logo-bw.svg',
      textLink: {
        text: 'ethermint.zone',
        url: 'https://ethermint.zone'
      },
      services: [
        {
          service: 'github',
          url: 'https://github.com/ChainSafe'
        },
        {
          service: 'twitter',
          url: 'https://twitter.com/chainsafeth'
        },
        {
          service: 'linkedin',
          url: 'https://www.linkedin.com/company/chainsafe-systems'
        },
        {
          service: 'medium',
          url: 'https://medium.com/chainsafe-systems'
        },
      ],
      smallprint:
          'This website is maintained by [ChainSafe Systems](https://chainsafe.io). The contents and opinions of this website are those of Chainsafe Systems.',
      links: [
        {
          title: 'Documentation',
          children: [
            {
              title: 'Cosmos SDK Docs',
              url: 'https://docs.cosmos.network'
            },
            {
              title: 'Ethermint Docs',
              url: 'https://ethereum.org/developers'
            },
            {
              title: 'Tendermint Core Docs',
              url: 'https://docs.tendermint.com'
            }
          ]
        },
        {
          title: 'Community',
          children: [
            {
              title: 'Cosmos Community',
              url: 'https://discord.gg/W8trcGV'
            },
            {
              title: 'Ethermint Forum',
              url: 'https://forum.cosmos.network/'
            },
            {
              title: 'Chainsafe Blog',
              url: 'https://medium.com/chainsafe-systems'
            }
          ]
        },
        {
          title: 'Contributing',
          children: [
            {
              title: 'Contributing to the docs',
              url:
                  'https://github.com/ChainSafe/ethermint/tree/development/docs'
            },
            {
              title: 'Careers at Chainsafe',
              url: 'https://chainsafe.io/#careers'
            },
            {
              title: 'Source code on GitHub',
              url: 'https://github.com/chainsafe/ethermint'
            }
          ]
        }
      ]
    }
  },
  plugins: [
    [
      '@vuepress/google-analytics',
      {
        ga: 'UA-51029217-12'
      }
    ],
    [
      'sitemap',
      {
        hostname: 'https://docs.cosmos.network'
      }
    ]
  ]
};
