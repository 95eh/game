package cmn

import "github.com/95eh/eg"

const (
	SceneNil eg.TScene = iota
	SceneWorld
)

func SceneTypeToSceneService(sceneType eg.TScene) (eg.TService, eg.IErr) {
	switch sceneType {
	case SceneWorld:
		return SvcSceneWorld, nil
	default:
		return SceneNil, eg.NewErr(eg.EcParamsErr, eg.M{
			"scene type": sceneType,
		})
	}
}
