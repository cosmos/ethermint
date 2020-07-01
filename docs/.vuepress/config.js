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
    repo: 'Chainsafe/ethermint',
    docsRepo: 'Chainsafe/ethermint',
    docsDir: 'docs',
    editLinks: true,
    label: 'ethermint',
    autoSidebar: false,
    sidebar: {
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
              path: '/modules',
              directory: true
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
      ],
    },
    gutter: {
      title: 'Help & Support',
      editLink: true,
      chat: {
        title: 'Discord',
        text: 'Chat with Cosmos developers on Discord.',
        url: 'https://discordapp.com/channels/669268347736686612',
        bg: 'linear-gradient(225.11deg, #2E3148 0%, #161931 95.68%)'
      },
      forum: {
        title: 'Cosmos Forum',
        text: 'Join the Cosmos Developer Forum to learn more.',
        url: 'https://forum.cosmos.network/',
        bg: 'linear-gradient(225deg, #46509F -1.08%, #2F3564 95.88%)',
        logo: 'cosmos'
      },
      github: {
        title: 'Found an Issue?',
        text: 'Help us improve this page by suggesting edits on GitHub.',
        url: 'https://github.com/Chainsafe/ethermint/edit/development/docs/README.md'  // FIXME: this is displayed to master
      }
    },
    footer: {
      questionsText: 'Chat with Cosmos developers on [Discord](https://discord.gg/W8trcGV) or reach out on the [SDK Developer Forum](https://forum.cosmos.network/) to learn more.',
      logo: '/logo-bw.svg',
      textLink: {
        text: 'cosmos.network',
        url: 'https://cosmos.network'
      },
      services: [
        {
          service: 'medium',
          url: 'https://blog.cosmos.network/'
        },
        {
          service: 'twitter',
          url: 'https://twitter.com/cosmos'
        },
        {
          service: 'github',
          url: 'https://github.com/ChainSafe/ethermint'
        },
      ],
      smallprint:
          'This website is maintained by Chainsafe Systems Inc. The contents and opinions of this website are those of Chainsafe Systems Inc.',
      links: [
        {
          title: 'Documentation',
          children: [
            {
              title: 'Cosmos SDK',
              url: 'https://docs.cosmos.network/'
            },
            {
              title: 'Ethermint',
              url: 'https://ethermint.zone/'
            },
            {
              title: 'Ethereum',
              url: 'https://ethereum.org/en/developers/'
            },
            {
              title: 'Tendermint Core',
              url: 'https://docs.tendermint.com/'
            }
          ]
        },
        {
          title: 'Community',
          children: [
            {
              title: 'Cosmos blog',
              url: 'https://blog.cosmos.network/'
            },
            {
              title: 'Forum',
              url: 'https://forum.cosmos.network/'
            },
            {
              title: 'Chat',
              url: 'https://riot.im/app/#/room/#cosmos-sdk:matrix.org'
            }
          ]
        },
        {
          title: 'Contributing',
          children: [
            {
              title: 'Contributing to the docs',
              url:
                  'https://github.com/Chainsafe/ethermint/blob/development/docs/DOCS_README.md'
            },
            {
              title: 'Source code on GitHub',
              url: 'https://github.com/Chainsafe/ethermint/'
            }
          ]
        }
      ]
    }
  },
};
