# API Reference

Packages:

- [tracing.orange.io/v1alpha1](#tracingorangeiov1alpha1)

# tracing.orange.io/v1alpha1

Resource Types:

- [Jaeger](#jaeger)




## Jaeger
<sup><sup>[↩ Parent](#tracingorangeiov1alpha1 )</sup></sup>






Jaeger is the Schema for the jaegers API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>tracing.orange.io/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Jaeger</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#jaegerspec">spec</a></b></td>
        <td>object</td>
        <td>
          JaegerSpec Define the desired state of Jaeger<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerstatus">status</a></b></td>
        <td>object</td>
        <td>
          JaegerStatus Define the observed state of Jaeger<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec
<sup><sup>[↩ Parent](#jaeger)</sup></sup>



JaegerSpec Define the desired state of Jaeger

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccommonspec">commonSpec</a></b></td>
        <td>object</td>
        <td>
          CommonSpec Define the configuration of components in Kubernetes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponents">components</a></b></td>
        <td>object</td>
        <td>
          Components Define the subComponents of Jaeger<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspecextensions">extensions</a></b></td>
        <td>object</td>
        <td>
          Extensions Define the configuration of the extension<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type Define the type of Jaeger deployment<br/>
          <br/>
            <i>Enum</i>: allInOne, distribute<br/>
            <i>Default</i>: allInOne<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.commonSpec
<sup><sup>[↩ Parent](#jaegerspec)</sup></sup>



CommonSpec Define the configuration of components in Kubernetes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccommonspecdeployment">deployment</a></b></td>
        <td>object</td>
        <td>
          Deployment Define configuration of the kubernetes Deployments<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccommonspecmetadata">metadata</a></b></td>
        <td>object</td>
        <td>
          Metadata Define metadata configuration of the component<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccommonspecservice">service</a></b></td>
        <td>object</td>
        <td>
          Service Define configuration of the kubernetes Services<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.commonSpec.deployment
<sup><sup>[↩ Parent](#jaegerspeccommonspec)</sup></sup>



Deployment Define configuration of the kubernetes Deployments

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>replicas</b></td>
        <td>integer</td>
        <td>
          Replicas Define Deployment replicas number<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>version</b></td>
        <td>string</td>
        <td>
          Version Define the version of the image used by the Jaeger instance component<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.commonSpec.metadata
<sup><sup>[↩ Parent](#jaegerspeccommonspec)</sup></sup>



Metadata Define metadata configuration of the component

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          Annotations Define annotations setting for metadata on the resource<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          Labels Define labels setting for metadata on the resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.commonSpec.service
<sup><sup>[↩ Parent](#jaegerspeccommonspec)</sup></sup>



Service Define configuration of the kubernetes Services

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Service Type string describes ingress methods for a service<br/>
          <br/>
            <i>Enum</i>: ClusterIP, NodePort, LoadBalancer<br/>
            <i>Default</i>: ClusterIP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components
<sup><sup>[↩ Parent](#jaegerspec)</sup></sup>



Components Define the subComponents of Jaeger

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccomponentsallinone">allInOne</a></b></td>
        <td>object</td>
        <td>
          AllInOne Define all-in-one component<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponentscollector">collector</a></b></td>
        <td>object</td>
        <td>
          Collector Define collector component<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponentsquery">query</a></b></td>
        <td>object</td>
        <td>
          QueryComponent Define collector component<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponentsstorage">storage</a></b></td>
        <td>object</td>
        <td>
          Storage Define backend storage component<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.allInOne
<sup><sup>[↩ Parent](#jaegerspeccomponents)</sup></sup>



AllInOne Define all-in-one component

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccomponentsallinonesetting">setting</a></b></td>
        <td>object</td>
        <td>
          ComponentSettings Define common Settings between components<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.allInOne.setting
<sup><sup>[↩ Parent](#jaegerspeccomponentsallinone)</sup></sup>



ComponentSettings Define common Settings between components

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Args Defined cmd args for Jaeger components<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponentsallinonesettingenvsindex">envs</a></b></td>
        <td>[]object</td>
        <td>
          Envs Defined env for Jaeger components<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.allInOne.setting.envs[index]
<sup><sup>[↩ Parent](#jaegerspeccomponentsallinonesetting)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Define Env name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Define Env value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.collector
<sup><sup>[↩ Parent](#jaegerspeccomponents)</sup></sup>



Collector Define collector component

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccomponentscollectorsetting">setting</a></b></td>
        <td>object</td>
        <td>
          ComponentSettings Define common Settings between components<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.collector.setting
<sup><sup>[↩ Parent](#jaegerspeccomponentscollector)</sup></sup>



ComponentSettings Define common Settings between components

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Args Defined cmd args for Jaeger components<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponentscollectorsettingenvsindex">envs</a></b></td>
        <td>[]object</td>
        <td>
          Envs Defined env for Jaeger components<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.collector.setting.envs[index]
<sup><sup>[↩ Parent](#jaegerspeccomponentscollectorsetting)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Define Env name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Define Env value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.query
<sup><sup>[↩ Parent](#jaegerspeccomponents)</sup></sup>



QueryComponent Define collector component

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccomponentsquerysetting">setting</a></b></td>
        <td>object</td>
        <td>
          ComponentSettings Define common Settings between components<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.query.setting
<sup><sup>[↩ Parent](#jaegerspeccomponentsquery)</sup></sup>



ComponentSettings Define common Settings between components

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Args Defined cmd args for Jaeger components<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspeccomponentsquerysettingenvsindex">envs</a></b></td>
        <td>[]object</td>
        <td>
          Envs Defined env for Jaeger components<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.query.setting.envs[index]
<sup><sup>[↩ Parent](#jaegerspeccomponentsquerysetting)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Define Env name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Define Env value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.storage
<sup><sup>[↩ Parent](#jaegerspeccomponents)</sup></sup>



Storage Define backend storage component

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspeccomponentsstoragees">es</a></b></td>
        <td>object</td>
        <td>
          Es Define backend storage is elasticsearch<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>options</b></td>
        <td>[]string</td>
        <td>
          Options Define backend storage options<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type Define backend storage type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.components.storage.es
<sup><sup>[↩ Parent](#jaegerspeccomponentsstorage)</sup></sup>



Es Define backend storage is elasticsearch

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.extensions
<sup><sup>[↩ Parent](#jaegerspec)</sup></sup>



Extensions Define the configuration of the extension

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#jaegerspecextensionshttproutesindex">httpRoutes</a></b></td>
        <td>[]object</td>
        <td>
          HTTPRoute Define the configuration of GatewayAPI routes<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.extensions.httpRoutes[index]
<sup><sup>[↩ Parent](#jaegerspecextensions)</sup></sup>



HTTPRoute Define the HTTPRoute configuration for the Gateway API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>hostnames</b></td>
        <td>[]string</td>
        <td>
          Hostnames Define a set of hostnames that should match against the HTTP Host
header to select a HTTPRoute used to process the request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#jaegerspecextensionshttproutesindexparentref">parentRef</a></b></td>
        <td>object</td>
        <td>
          ParentRefs references the resources (usually Gateways) that a Route wants
to be attached to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>target</b></td>
        <td>enum</td>
        <td>
          Target Define the Service target to which routes need to be added<br/>
          <br/>
            <i>Enum</i>: collector, query<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetPort</b></td>
        <td>integer</td>
        <td>
          TargetPort Define the Service target port to which routes need to be added<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.spec.extensions.httpRoutes[index].parentRef
<sup><sup>[↩ Parent](#jaegerspecextensionshttproutesindex)</sup></sup>



ParentRefs references the resources (usually Gateways) that a Route wants
to be attached to.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of the referent.


Support: Core<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>group</b></td>
        <td>string</td>
        <td>
          Group is the group of the referent.
When unspecified, "gateway.networking.k8s.io" is inferred.
To set the core API group (such as for a "Service" kind referent),
Group must be explicitly set to "" (empty string).


Support: Core<br/>
          <br/>
            <i>Default</i>: gateway.networking.k8s.io<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is kind of the referent.


There are two kinds of parent resources with "Core" support:


* Gateway (Gateway conformance profile)
* Service (Mesh conformance profile, experimental, ClusterIP Services only)


Support for other resources is Implementation-Specific.<br/>
          <br/>
            <i>Default</i>: Gateway<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace is the namespace of the referent. When unspecified, this refers
to the local namespace of the Route.


Note that there are specific rules for ParentRefs which cross namespace
boundaries. Cross-namespace references are only valid if they are explicitly
allowed by something in the namespace they are referring to. For example:
Gateway has the AllowedRoutes field, and ReferenceGrant provides a
generic way to enable any other kind of cross-namespace reference.


<gateway:experimental:description>
ParentRefs from a Route to a Service in the same namespace are "producer"
routes, which apply default routing rules to inbound connections from
any namespace to the Service.


ParentRefs from a Route to a Service in a different namespace are
"consumer" routes, and these routing rules are only applied to outbound
connections originating from the same namespace as the Route, for which
the intended destination of the connections are a Service targeted as a
ParentRef of the Route.
</gateway:experimental:description>


Support: Core<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port is the network port this Route targets. It can be interpreted
differently based on the type of parent resource.


When the parent resource is a Gateway, this targets all listeners
listening on the specified port that also support this kind of Route(and
select this Route). It's not recommended to set `Port` unless the
networking behaviors specified in a Route must apply to a specific port
as opposed to a listener(s) whose port(s) may be changed. When both Port
and SectionName are specified, the name and port of the selected listener
must match both specified values.


<gateway:experimental:description>
When the parent resource is a Service, this targets a specific port in the
Service spec. When both Port (experimental) and SectionName are specified,
the name and port of the selected port must match both specified values.
</gateway:experimental:description>


Implementations MAY choose to support other parent resources.
Implementations supporting other types of parent resources MUST clearly
document how/if Port is interpreted.


For the purpose of status, an attachment is considered successful as
long as the parent resource accepts it partially. For example, Gateway
listeners can restrict which Routes can attach to them by Route kind,
namespace, or hostname. If 1 of 2 Gateway listeners accept attachment
from the referencing Route, the Route MUST be considered successfully
attached. If no Gateway listeners accept attachment from this Route,
the Route MUST be considered detached from the Gateway.


Support: Extended


<gateway:experimental><br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Minimum</i>: 1<br/>
            <i>Maximum</i>: 65535<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sectionName</b></td>
        <td>string</td>
        <td>
          SectionName is the name of a section within the target resource. In the
following resources, SectionName is interpreted as the following:


* Gateway: Listener Name. When both Port (experimental) and SectionName
are specified, the name and port of the selected listener must match
both specified values.
* Service: Port Name. When both Port (experimental) and SectionName
are specified, the name and port of the selected listener must match
both specified values. Note that attaching Routes to Services as Parents
is part of experimental Mesh support and is not supported for any other
purpose.


Implementations MAY choose to support attaching Routes to other resources.
If that is the case, they MUST clearly document how SectionName is
interpreted.


When unspecified (empty string), this will reference the entire resource.
For the purpose of status, an attachment is considered successful if at
least one section in the parent resource accepts it. For example, Gateway
listeners can restrict which Routes can attach to them by Route kind,
namespace, or hostname. If 1 of 2 Gateway listeners accept attachment from
the referencing Route, the Route MUST be considered successfully
attached. If no Gateway listeners accept attachment from this Route, the
Route MUST be considered detached from the Gateway.


Support: Core<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.status
<sup><sup>[↩ Parent](#jaeger)</sup></sup>



JaegerStatus Define the observed state of Jaeger

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>phase</b></td>
        <td>string</td>
        <td>
          Phase Define the component phase of Jaeger<br/>
          <br/>
            <i>Default</i>: Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#jaegerstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions  Define the conditions of Jaeger<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Jaeger.status.conditions[index]
<sup><sup>[↩ Parent](#jaegerstatus)</sup></sup>



Condition contains details for one aspect of the current state of this API Resource.
---
This struct is intended for direct use as an array at the field path .status.conditions.  For example,


	type FooStatus struct{
	    // Represents the observations of a foo's current state.
	    // Known .status.conditions.type are: "Available", "Progressing", and "Degraded"
	    // +patchMergeKey=type
	    // +patchStrategy=merge
	    // +listType=map
	    // +listMapKey=type
	    Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`


	    // other fields
	}

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          lastTransitionTime is the last time the condition transitioned from one status to another.
This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.
This may be an empty string.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.
Producers of specific condition types may define expected values and meanings for this field,
and whether the values are considered a guaranteed API.
The value should be a CamelCase string.
This field may not be empty.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          status of the condition, one of True, False, Unknown.<br/>
          <br/>
            <i>Enum</i>: True, False, Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type of condition in CamelCase or in foo.example.com/CamelCase.
---
Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be
useful (see .node.status.conditions), the ability to deconflict is important.
The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
