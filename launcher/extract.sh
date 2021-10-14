#!/bin/bash

NAMESPACE_TESTRUNNERPODS=$(cat ${1} | jq -r '.testRunnerPods[] | .namespace + "/" + .podName')

echo '############### extract data ###############'

mkdir testresults

for NAMESPACE_TESTRUNNERPOD in ${NAMESPACE_TESTRUNNERPODS} ; do

    NAMESPACE=$(dirname ${NAMESPACE_TESTRUNNERPOD})
    PODNAME=$(basename ${NAMESPACE_TESTRUNNERPOD})

    SIDECAR_CONTAINERS=$(cat ${1} | jq -r '.testRunnerPods[] | select(.podName=='\"$PODNAME\"') | .testRunnerSidecarContainers[] | .sidecarContainerName')

    for SIDECAR_CONTAINER in "${SIDECAR_CONTAINERS[@]}" ; do

        kubectl --namespace=${NAMESPACE} cp "${NAMESPACE_TESTRUNNERPOD}:/testresults" "testresults/${NAMESPACE}-${PODNAME}" -c "${SIDECAR_CONTAINER}"

    done

done

echo '############### delete test runner pods ###############'

for row in ${NAMESPACE_TESTRUNNERPODS} ; do

    NAMESPACE=$(dirname ${row})
    PODNAME=$(basename ${row})

    kubectl --namespace=${NAMESPACE} delete pod ${PODNAME}

done