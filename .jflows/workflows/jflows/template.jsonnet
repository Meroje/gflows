local git_config = import '../config/git.libsonnet';

local check_workflows_job = {
  'name': 'check-workflows',
  'runs-on': 'ubuntu-latest',
  steps: [
    {
      uses: 'actions/checkout@v2'
    },
    {
      name: 'install jflows',
      run: 'go get github.com/jbrunton/jflows'
    },
    {
      name: 'validate workflows',
      run: 'jflows check'
    },
  ]
};

local workflow = {
  name: 'g3ops-workflows',
  on: {
    pull_request: {
      branches: [git_config.main_branch]
    },
    push: {
      branches: [git_config.main_branch]
    }
  },
  jobs: {
    check_workflows: check_workflows_job
  },
};

std.manifestYamlDoc(workflow)